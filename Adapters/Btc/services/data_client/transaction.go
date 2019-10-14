package data_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/services/converter"
	"net/http"
)

const (
	txIsNotInBlockchain = "Cannot find transaction with id"
	sendTxUrl           = "/transactions"
	unconfirmedTxUrl    = "/transactions/unconfirmed"
	TxByIdUrlTemplate   = "/transactions/%s"
	getTxByHashUrl      = "/tx/%s"
)


func (dcl *dataClient) TransactionByHash(ctx context.Context, txId string) (*models.TxInfo, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'TransactionByHash' of dataClient for txID %s", txId)

	txResp, err := dcl.Request(ctx, http.MethodGet, dcl.conf.Url+fmt.Sprintf(getTxByHashUrl, txId), nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	tx := &models.RawTx{}
	if err := json.Unmarshal(txResp, tx); err != nil {
		log.Errorf("failed to unmarshal raw tx: %s", err)
		return nil, err
	}

	return dcl.parseTx(tx), nil
}


func (dcl *dataClient) parseTx(tx *models.RawTx) *models.TxInfo {
	inputs := make([]models.InputOutputInfo, 0)
	outputs := make([]models.InputOutputInfo, 0)

	for _, input := range tx.Inputs {
		inputs = append(inputs, models.InputOutputInfo{
			Amount:  converter.ToTargetAmountStr(input.Value),
			Address: input.Address,
		})
	}

	amount := tx.ValueIn - tx.ValueOut
	for _, output := range tx.Outputs {
		outputs = append(outputs, models.InputOutputInfo{
			Amount: output.Value,
		})
	}
	return &models.TxInfo{
		Amount:  converter.ToTargetAmountStr(amount),
		TxHash:  tx.Id,
		Status:  models.TxStatusSuccess,
		Inputs:  inputs,
		Outputs: outputs,
		Fee:     converter.ToTargetAmountStr(tx.Fees),
	}
}
