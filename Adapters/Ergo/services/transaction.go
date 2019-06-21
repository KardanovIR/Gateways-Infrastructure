package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/models"
)

const (
	txIsNotInBlockchain = "Transaction is not in blockchain"
	sendTxUrl           = "/transactions"
)

type SendTxResponse struct {
	ID string `json:"id"`
}

func (cl *nodeClient) SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'SendTransaction'")
	sendTxResp, _ := cl.Request(ctx, http.MethodPost, cl.conf.ExplorerUrl+sendTxUrl, bytes.NewReader(txSigned))
	sendTxResponse := SendTxResponse{}
	if err := json.Unmarshal(sendTxResp, &sendTxResponse); err != nil {
		log.Errorf("failed to send tx %s: %s", string(txSigned), err)
		return "", err
	}
	// explorer returns tx id with quotes - replace them
	txId = replaceQuotesFromSides(sendTxResponse.ID)
	return txId, nil
}

func replaceQuotesFromSides(s string) string {
	return strings.Replace(s, "\"", "", 2)
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
