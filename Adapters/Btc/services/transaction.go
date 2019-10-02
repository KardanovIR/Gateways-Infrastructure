package services

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

type SendTxResponse struct {
	ID string `json:"id"`
}

func (cl *nodeClient) SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'SendRawTransaction'")
	//todo
	//txHash, err := cl.nodeClient.SendRawTransaction(txSigned, false)
	//if err != nil {
	//	log.Error(err)
	//	return "", err
	//}

	//log.Debugf("node return %s", txHash)
	//return *txHash, nil
	return "", nil
}

func (cl *dataClient) TransactionByHash(ctx context.Context, txId string) (*models.TxInfo, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'TransactionByHash' for txID %s", txId)

	txResp, err := cl.Request(ctx, http.MethodGet, cl.conf.Url+fmt.Sprintf(getTxByHashUrl, txId), nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	tx := &models.RawTx{}
	if err := json.Unmarshal(txResp, tx); err != nil {
		log.Errorf("failed to unmarshal raw tx: %s", err)
		return nil, err
	}

	return parseTx(tx), nil
}

func parseTx(tx *models.RawTx) *models.TxInfo {
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

func (cl *nodeClient) Fee(ctx context.Context, senderPublicKey string) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'Fee'")
	//todo
	return 0, nil
}
