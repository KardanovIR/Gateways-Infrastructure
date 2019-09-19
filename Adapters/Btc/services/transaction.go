package services

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
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

func (cl *nodeClient) TransactionByHash(ctx context.Context, txId string) (*models.TxInfo, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'TransactionByHash' for txID %s", txId)
	txHash,err  := chainhash.NewHashFromStr(txId)
	log.Infof("new hash %s", txHash)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Infof("call node")
	nodeTx, err := cl.nodeClient.GetRawTransactionVerbose(txHash)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Infof("node's response %s", nodeTx)
	return parseTx(nodeTx), nil
}

func parseTx(tx *btcjson.TxRawResult) *models.TxInfo {

	inputs := make([]models.InputOutputInfo, 0)
	outputs := make([]models.InputOutputInfo, 0)

	for _, input := range tx.Vin {
		//todo доделать
		inputs = append(inputs, models.InputOutputInfo{
			//Amount: fmt.Sprintf("%f", input),
			Address: input.Txid,
		})
	}

	amount:= 0.0
	for _, output := range tx.Vout {
		if len(output.ScriptPubKey.Addresses) == 0 {
			continue
		}
		inputs = append(inputs, models.InputOutputInfo{
			Amount: fmt.Sprintf("%f", output.Value),
			Address: output.ScriptPubKey.Addresses[0],
		})
		amount += output.Value
	}

	return &models.TxInfo{
		Amount:  fmt.Sprintf("%f", amount),
		TxHash:  tx.Txid,
		Status:  models.TxStatusSuccess,
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func (cl *nodeClient) Fee(ctx context.Context, senderPublicKey string) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'Fee'")
	//todo
	return 0, nil
}