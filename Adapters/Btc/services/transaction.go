package services

import (
	"context"
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
	var chainHash *chainhash.Hash
	err := chainhash.Decode(chainHash, txId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	_, err = cl.nodeClient.GetTransaction(chainHash)
	//&models.TxInfo{
	//	TxHash:  nodeTx.Hex,
	//	Amount:  nodeTx.Amount,
	//	Fee: nodeTx.Fee,
	//
	//	Inputs:  inputs,
	//	Outputs: outputs,
	//}
	//nodeTx.Details

	return nil, nil
}

func parseTx(tx *btcjson.GetTransactionResult) *models.TxInfo {
	//todo переделать под btc
	//inputs := make([]models.InputOutputInfo, 0)
	//outputs := make([]models.InputOutputInfo, 0)



	return &models.TxInfo{
		//From:    sender,
		//To:      recipient,
		//Fee:     converter.ToTargetAmountStr(fee),
		//Amount:  amount,
		//TxHash:  tx.Summary.ID,
		//Status:  models.TxStatusSuccess,
		//Inputs:  inputs,
		//Outputs: outputs,
	}
}

func (cl *nodeClient) Fee(ctx context.Context, senderPublicKey string) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'Fee'")
	//todo
	return 0, nil
}