package services

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"golang.org/x/crypto/sha3"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/models"
)

var (
	transferErc20MethodID []byte
	initTransferOnce      sync.Once
)

func (cl *nodeClient) CreateRawTransaction(ctx context.Context, addressFrom string, addressTo string,
	amount *big.Int) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'CreateRawTransaction': send %s from %s to %s", amount, addressFrom, addressTo)
	gasPrice, err := cl.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get suggected gas price %s", err)
	}
	log.Debugf("suggest gas price %s", gasPrice)
	nonce, err := cl.GetNextNonce(ctx, addressFrom)
	if err != nil {
		return nil, fmt.Errorf("can't get next nonce for address %s: %s", addressFrom, err)
	}
	log.Debugf("nonce will be %d", nonce)
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(addressTo),
		amount,
		gasLimitForMoneyTransfer,
		gasPrice,
		nil,
	)
	b, err := SerializeTx(tx)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (cl *nodeClient) CreateErc20TokensRawTransaction(ctx context.Context, addressFrom string, contractAddress string,
	addressTo string, amount *big.Int) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'CreateErc20TokensRawTransaction': send %s tokens from %s to %s; contract %s",
		amount, addressFrom, addressTo, contractAddress)

	gasPrice, err := cl.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get suggected gas price %s", err)
	}
	log.Debugf("suggest gas price %s", gasPrice)
	nonce, err := cl.GetNextNonce(ctx, addressFrom)
	log.Debugf("nonce will be %d", nonce)
	if err != nil {
		return nil, fmt.Errorf("can't get next nonce for address %s: %s", addressFrom, err)
	}
	recipient := common.HexToAddress(addressTo)
	sender := common.HexToAddress(addressFrom)
	tokenAddress := common.HexToAddress(contractAddress)
	methodID := getTransferErc20MethodID()

	// add zeros before address and amount value
	paddedAddress := common.LeftPadBytes(recipient.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := cl.ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		From: sender,
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("estimated gas limit %d", gasLimit)
	ethAmountZero := big.NewInt(0)
	tx := types.NewTransaction(nonce, tokenAddress, ethAmountZero, gasLimit, gasPrice, data)
	b, err := SerializeTx(tx)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (cl *nodeClient) SignTransaction(ctx context.Context, senderAddr string, rlpTx []byte) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'SignTransaction'")
	pk, ok := cl.privateKeys[senderAddr]
	if !ok {
		return nil, fmt.Errorf("can't signed transaction for %s: haven't private key", senderAddr)
	}
	return cl.signTx(ctx, pk, rlpTx)
}

func (cl *nodeClient) SignTransactionWithPrivateKey(ctx context.Context, privateKey string, rlpTx []byte) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'SignTransactionWithPrivateKey'")
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Errorf("can't cast key to ECDSA: %s", err)
		return nil, err
	}
	return cl.signTx(ctx, pk, rlpTx)
}

func (cl *nodeClient) signTx(ctx context.Context, privateKey *ecdsa.PrivateKey, rlpTx []byte) ([]byte, error) {
	tx, err := DeserializeTx(rlpTx)
	if err != nil {
		return nil, fmt.Errorf("can't deserialize tx %s", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(cl.chainID)), privateKey)
	if err != nil {
		return nil, err
	}
	return SerializeTx(signedTx)
}

func (cl *nodeClient) SendTransaction(ctx context.Context, rlpTx []byte) (txHex string, err error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'SendTransaction'")

	tx, err := DeserializeTx(rlpTx)
	if err != nil {
		return "", fmt.Errorf("can't deserialize tx %s", err)
	}
	log.Debugf("try to send transaction %+v", tx)
	return tx.Hash().Hex(), cl.ethClient.SendTransaction(ctx, tx)
}

func (cl *nodeClient) GetTxStatusByTxID(ctx context.Context, txID string) (models.TxStatus, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'GetTransactionByHash' %s", txID)
	_, pending, err := cl.ethClient.TransactionByHash(ctx, common.HexToHash(txID))
	if err != nil {
		if err == ethereum.NotFound {
			return models.TxStatusUnKnown, nil
		}
		return "", err
	}
	if pending {
		return models.TxStatusPending, nil
	}
	return models.TxStatusSuccess, nil
}

func DeserializeTx(rlpTx []byte) (*types.Transaction, error) {
	reader := bytes.NewReader(rlpTx)
	tx := new(types.Transaction)
	rlpStream := rlp.NewStream(reader, 0)
	if err := tx.DecodeRLP(rlpStream); err != nil {
		return nil, err
	}
	return tx, nil
}

func SerializeTx(tx *types.Transaction) ([]byte, error) {
	var b bytes.Buffer
	if err := tx.EncodeRLP(&b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func getTransferErc20MethodID() []byte {
	initTransferOnce.Do(func() {
		transferFnSignature := []byte("transfer(address,uint256)")
		hash := sha3.NewLegacyKeccak256()
		hash.Write(transferFnSignature)
		transferErc20MethodID = hash.Sum(nil)[:4]
	})
	return transferErc20MethodID
}
