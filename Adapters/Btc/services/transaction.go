package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/_vendor-20180717124023/github.com/btcsuite/btcd/chaincfg/chainhash"
	"net/http"
	"strings"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/services/converter"
)

const (
	txIsNotInBlockchain = "Cannot find transaction with id"
	sendTxUrl           = "/transactions"
	unconfirmedTxUrl    = "/transactions/unconfirmed"
	TxByIdUrlTemplate   = "/transactions/%s"
)

type SendTxResponse struct {
	ID string `json:"id"`
}

func (cl *nodeClient) SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'SendRawTransaction'")
	txHash, err := cl.nodeClient.SendRawTransaction(bytes.NewReader(txSigned), false)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Debugf("node return %s", txHash)
	return *txHash, nil
}

func (cl *nodeClient) TransactionByHash(ctx context.Context, txId string) (*models.TxInfo, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'TransactionByHash' for txID %s", txId)
	var chainHash *chainhash.Hash
	err := chainhash.Decode(chainHash, txId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	nodeTx, err := cl.nodeClient.GetTransaction(chainHash)
	_ := nodeTx

	return parseTx(nodeTx), nil
}

func parseTx(tx *models.Tx) *models.TxInfo {
	//todo переделать под btc
	inputs := make([]models.InputOutputInfo, 0)
	outputs := make([]models.InputOutputInfo, 0)

	return &models.TxInfo{
		From:    sender,
		To:      recipient,
		Fee:     converter.ToTargetAmountStr(fee),
		Amount:  amount,
		TxHash:  tx.Summary.ID,
		Status:  models.TxStatusSuccess,
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func (cl *nodeClient) Fee(ctx context.Context, senderPublicKey string) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'Fee'")
	//todo
	log.Debugf("node return %s", txHash)
	return 0, nil
}