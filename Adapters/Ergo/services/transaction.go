package services

import (
	"context"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/models"
)

const (
	txIsNotInBlockchain = "Transaction is not in blockchain"
)

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
