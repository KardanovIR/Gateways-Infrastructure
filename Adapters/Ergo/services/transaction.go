package services

import (
	"context"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/models"
)

const (
	decimalBase         = 10
	txIsNotInBlockchain = "Transaction is not in blockchain"
)

// CreateRawTxBySendersAddress creates transaction for senders address if private key keeps in adapter
func (cl *nodeClient) CreateRawTxBySendersAddress(ctx context.Context, addressFrom string,
	addressTo string, amount uint64) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'CreateRawTxBySendersAddress' from %s to %s amount %d",
		addressFrom, addressTo, amount)
	// todo implementation
	return []byte{}, nil
}

// CreateRawTxBySendersPublicKey creates transaction using public key. Private key is not used
func (cl *nodeClient) CreateRawTxBySendersPublicKey(ctx context.Context, sendersPublicKey string,
	addressTo string, amount uint64, assetId string) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'CreateRawTxBySendersPublicKey' pk %s send amount %d to %s",
		sendersPublicKey, amount, addressTo)
	// todo implementation
	return []byte{}, nil
}

func (cl *nodeClient) SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'SendTransaction'")
	// todo implementation
	return
}

func (cl *nodeClient) GetTransactionByTxId(ctx context.Context, txId string) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'GetTransactionByTxId' for txID %s", txId)
	// todo implementation
	return []byte{}, nil
}

func (cl *nodeClient) TransactionByHash(ctx context.Context, txId string) (*models.TxInfo, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'TransactionByHash' for txID %s", txId)
	// todo implementation
	return &models.TxInfo{}, nil
}

func (cl *nodeClient) GetTransactionStatus(ctx context.Context, txId string) (models.TxStatus, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'GetTransactionStatus' for txID %s", txId)
	// todo implementation
	return models.TxStatusSuccess, nil
}
