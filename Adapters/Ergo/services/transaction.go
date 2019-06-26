package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/services/converter"
)

const (
	txIsNotInBlockchain = "Cannot find transaction with id"
	sendTxUrl           = "/transactions"
	unconfirmedTxUrl    = "/transactions/unconfirmed"
	TxByIdUrlTemplate   = "/transactions/%s"
	decimalBase         = 10
)

type SendTxResponse struct {
	ID string `json:"id"`
}

func (cl *nodeClient) SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'SendTransaction'")
	sendTxResp, err := cl.Request(ctx, http.MethodPost, cl.conf.ExplorerUrl+sendTxUrl, bytes.NewReader(txSigned))
	if err != nil {
		log.Error(err)
		return "", err
	}
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

func (cl *nodeClient) TransactionByHash(ctx context.Context, txId string) (*models.TxInfo, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'TransactionByHash' for txID %s", txId)
	unconfirmedTxResp, err := cl.Request(ctx, http.MethodGet, cl.conf.ExplorerUrl+unconfirmedTxUrl, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	unconfirmedTxList := make([]*models.UnSignedTx, 0)
	if err := json.Unmarshal(unconfirmedTxResp, &unconfirmedTxList); err != nil {
		log.Errorf("failed to unmarshal unconfirmed tx list: %s", err)
		return nil, err
	}
	unconfirmedTx := findTxInList(unconfirmedTxList, txId)
	if unconfirmedTx != nil {
		// don't parse all info: not used by another service if tx is unconfirmed. But in fact it should be parsed
		return &models.TxInfo{Status: models.TxStatusPending}, nil
	}
	txResp, err := cl.Request(ctx, http.MethodGet, cl.conf.ExplorerUrl+fmt.Sprintf(TxByIdUrlTemplate, txId), nil)
	if err != nil {
		if e, ok := err.(*WrongCodeError); ok {
			if e.Code == 404 && strings.Contains(e.Body, txIsNotInBlockchain) {
				return &models.TxInfo{Status: models.TxStatusUnKnown}, nil
			}
		}
		log.Error(err)
		return nil, err
	}
	tx := models.Tx{}
	if err := json.Unmarshal(txResp, &tx); err != nil {
		log.Errorf("failed to unmarshal tx: %s", err)
		return nil, err
	}
	info := parseTx(&tx)
	return info, nil
}

func parseTx(tx *models.Tx) *models.TxInfo {
	inputs := make([]models.InputOutputInfo, 0)
	outputs := make([]models.InputOutputInfo, 0)
	feeId := findFeeOutput(tx.Outputs)
	fee := uint64(0)
	if feeId != nil {
		fee = tx.Outputs[*feeId].Value
		// delete fee output from outputs
		tx.Outputs[*feeId] = tx.Outputs[len(tx.Outputs)-1]
		tx.Outputs = tx.Outputs[:len(tx.Outputs)-1]
	}
	txInputs := summarizeAmountByAddress(tx.Inputs)
	txOutputs := summarizeAmountByAddress(tx.Outputs)
	for _, in := range txInputs {
		hasOut := false
		for _, out := range txOutputs {
			if in.Address == out.Address {
				hasOut = true
				if in.Value > out.Value {
					inputs = append(inputs, models.InputOutputInfo{
						Amount:  strconv.FormatUint(converter.ToTargetAmount(in.Value-out.Value), decimalBase),
						Address: in.Address,
					})
				} else if in.Value < out.Value {
					outputs = append(outputs, models.InputOutputInfo{
						Amount:  strconv.FormatUint(converter.ToTargetAmount(out.Value-in.Value), decimalBase),
						Address: out.Address,
					})
				}
				break
			}
		}
		if !hasOut {
			inputs = append(inputs, models.InputOutputInfo{
				Amount:  strconv.FormatUint(converter.ToTargetAmount(in.Value), decimalBase),
				Address: in.Address,
			})
		}
	}
	for _, out := range txOutputs {
		if !hasAddress(out.Address, txInputs) {
			outputs = append(outputs, models.InputOutputInfo{
				Amount:  strconv.FormatUint(converter.ToTargetAmount(out.Value), decimalBase),
				Address: out.Address,
			})
		}
	}
	sender := ""
	if len(inputs) == 1 {
		sender = inputs[0].Address
	}
	recipient := ""
	amount := ""
	if len(outputs) == 1 {
		recipient = outputs[0].Address
		amount = outputs[0].Amount
	}

	return &models.TxInfo{
		From:    sender,
		To:      recipient,
		Fee:     strconv.FormatUint(converter.ToTargetAmount(fee), decimalBase),
		Amount:  amount,
		TxHash:  tx.Summary.ID,
		Status:  models.TxStatusSuccess,
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func summarizeAmountByAddress(list []models.InputOutput) []models.InputOutput {
	result := make([]models.InputOutput, 0)
	for i := 0; i < len(list); i++ {
		a := list[i]
		if hasAddress(a.Address, result) {
			continue
		}
		amount := a.Value
		for j := i + 1; j < len(list); j++ {
			next := list[j]
			if next.Address == a.Address {
				amount += next.Value
			}
		}
		result = append(result, models.InputOutput{Address: a.Address, Value: amount, ErgoTree: a.ErgoTree})
	}
	return result
}

func findFeeOutput(list []models.InputOutput) *int {
	for i, a := range list {
		if a.ErgoTree == minerErgoTree {
			return &i
		}
	}
	return nil
}

func hasAddress(address string, list []models.InputOutput) bool {
	for _, t := range list {
		if address == t.Address {
			return true
		}
	}
	return false
}

func findTxInList(txList []*models.UnSignedTx, txID string) *models.UnSignedTx {
	for _, unTx := range txList {
		if txID == unTx.ID {
			return unTx
		}
	}
	return nil
}
