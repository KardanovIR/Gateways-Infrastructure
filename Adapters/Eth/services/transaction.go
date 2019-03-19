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
	amount *big.Int) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'CreateRawTransaction': send %s from %s to %s", amount, addressFrom, addressTo)
	gasPrice, err := cl.SuggestGasPrice(ctx)
	log.Debugf("suggest gas price %s", gasPrice)
	nonce, err := cl.GetNextNonce(ctx, addressFrom)
	log.Debugf("nonce will be %d", nonce)
	if err != nil {
		return nil, fmt.Errorf("can't get suggected gas price %s", err)
	}
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
