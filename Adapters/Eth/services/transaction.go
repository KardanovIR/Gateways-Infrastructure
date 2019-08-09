package services

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/models"
)

func (cl *nodeClient) CreateRawTransaction(ctx context.Context, addressFrom string, addressTo string,
	amount *big.Int, nonce uint64) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'CreateRawTransaction': send %s from %s to %s", amount, addressFrom, addressTo)
	ok, _, err := cl.IsAddressValid(ctx, addressTo)
	if err != nil {
		return nil, fmt.Errorf("check address %s fails: %s", addressTo, err)
	}
	if !ok {
		return nil, fmt.Errorf("address %s is not valid", addressTo)
	}
	gasPrice, err := cl.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get suggected gas price %s", err)
	}
	log.Debugf("suggest gas price %s", gasPrice)
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
	addressTo string, amount *big.Int, nonce uint64) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'CreateErc20TokensRawTransaction': send %s tokens from %s to %s; contract %s",
		amount, addressFrom, addressTo, contractAddress)
	ok, _, err := cl.IsAddressValid(ctx, addressTo)
	if err != nil {
		return nil, fmt.Errorf("check address %s fails: %s", addressTo, err)
	}
	if !ok {
		return nil, fmt.Errorf("address %s is not valid", addressTo)
	}
	gasPrice, err := cl.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get suggected gas price %s", err)
	}
	log.Debugf("suggest gas price %s", gasPrice)
	sender := common.HexToAddress(addressFrom)
	tokenAddress := common.HexToAddress(contractAddress)
	data := cl.contractProvider.CreateTransferTokenData(ctx, addressTo, amount)

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
	_, status, err := cl.getTxAndStatus(ctx, txID)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (cl *nodeClient) TransactionInfo(ctx context.Context, txID string) (*models.TxInfo, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'TransactionInfo' %s", txID)
	tx, status, err := cl.getTxAndStatus(ctx, txID)
	if err != nil {
		return nil, err
	}
	if status == models.TxStatusUnKnown {
		return &models.TxInfo{Status: status}, nil
	}
	sender, err := types.Sender(types.NewEIP155Signer(big.NewInt(cl.chainID)), tx)
	if err != nil {
		log.Errorf("can't get sender for tx %s: %s", txID, err)
		return nil, err
	}
	// Todo do for erc-20 tokens: another implementation for recipient, amount
	fee := new(big.Int).Mul(big.NewInt(int64(tx.Gas())), tx.GasPrice())
	txInfo := models.TxInfo{
		From:     sender.String(),
		To:       tx.To().String(),
		Amount:   tx.Value(),
		Fee:      fee,
		Contract: "",
		TxHash:   tx.Hash().String(),
		Data:     string(tx.Data()),
		Status:   status,
	}
	return &txInfo, nil
}

func (cl *nodeClient) getTxAndStatus(ctx context.Context, txID string) (*types.Transaction, models.TxStatus, error) {
	tx, pending, err := cl.ethClient.TransactionByHash(ctx, common.HexToHash(txID))
	if err != nil {
		if err == ethereum.NotFound {
			return &types.Transaction{}, models.TxStatusUnKnown, nil
		}
		return nil, "", err
	}
	if pending {
		return tx, models.TxStatusPending, nil
	}
	return tx, models.TxStatusSuccess, nil
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
