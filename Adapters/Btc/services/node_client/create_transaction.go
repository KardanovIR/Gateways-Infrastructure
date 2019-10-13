package node_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
)

const (
	MinBtcOutputValueConst = uint64(1000)
	// if locktime is 0 -> tx will be immediately processing (if locktime > 0 -> node waits for block or time before takes it)
	WithoutLockTime = int64(0)
)

type TxRawInfo struct {
	InputsInfo []TxInputInfo `json:"inputsInfo"`
	RawTx      []byte        `json:"rawTx"`
}

type TxInputInfo struct {
	Address string `json:"address"`
	Input   btcjson.TransactionInput
}

// CreateRawTx creates transaction
func (cl *nodeClient) CreateRawTx(ctx context.Context, addressesFrom []string, changeAddress string,
	outs []*models.Output) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'CreateRawTx' from %s to %v",
		addressesFrom, outs)
	if len(addressesFrom) == 0 || len(outs) == 0 {
		return nil, fmt.Errorf("wrong parameters addressesFrom %s or outs", addressesFrom)
	}
	amount := uint64(0)
	for _, o := range outs {
		if o.Amount < MinBtcOutputValueConst {
			return nil, fmt.Errorf("amount %d to sent to %s is less than min amount of fee", o.Amount, o.Address)
		}
		amount += o.Amount
	}
	fee, err := cl.Fee(ctx)
	if err != nil {
		log.Errorf("get fee fails %s", err)
		return nil, err
	}
	inputInfos, summaryInputsAmount, err := cl.getUnspentInputs(ctx, addressesFrom, amount+fee)
	if err != nil {
		log.Errorf("get UnspentInputs fails %s", err)
		return nil, err
	}
	inputsForTx := make([]btcjson.TransactionInput, 0)
	for _, inputInfo := range inputInfos {
		inputsForTx = append(inputsForTx, inputInfo.Input)
	}
	destinations := make(map[btcutil.Address]btcutil.Amount)
	change := summaryInputsAmount - (amount + fee)
	if change > MinBtcOutputValueConst {
		sendersAddress, err := btcutil.DecodeAddress(changeAddress, cl.conf.ChainParams)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		destinations[sendersAddress] = btcutil.Amount(change)
	}
	for _, out := range outs {
		outAddress, err := btcutil.DecodeAddress(out.Address, cl.conf.ChainParams)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		destinations[outAddress] = btcutil.Amount(out.Amount)
	}
	lockTime := WithoutLockTime
	rawTx, err := cl.nodeClient.CreateRawTransaction(inputsForTx, destinations, &lockTime)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	rt, err := Serialize(rawTx)
	// put info about address from which input will be get because we can't get address from txInput
	txRawInfo := TxRawInfo{RawTx: rt, InputsInfo: inputInfos}
	txBytesResp, err := json.Marshal(&txRawInfo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return txBytesResp, err
}

func (cl *nodeClient) getUnspentInputs(ctx context.Context, addressesFrom []string, summaryAmount uint64) ([]TxInputInfo, uint64, error) {
	log := logger.FromContext(ctx)
	unspentInputs, err := cl.rep.GetUnspentTxListForAddress(ctx, addressesFrom)
	if err != nil {
		log.Errorf("get unspent inputs fails %s", err)
		return nil, 0, err
	}
	summaryInputsAmount := uint64(0)
	inputInfos := make([]TxInputInfo, 0)
	// todo more efficient algorithm of choosing inputs
	for _, input := range unspentInputs {
		if summaryInputsAmount >= summaryAmount {
			break
		}
		txInput := btcjson.TransactionInput{Txid: input.TxHash, Vout: input.OutputN}
		inputInfos = append(inputInfos, TxInputInfo{Address: input.Address, Input: txInput})
		summaryInputsAmount += input.Amount
	}
	if summaryInputsAmount < summaryAmount {
		return nil, 0, fmt.Errorf("insufficient funds: need %d, has %d", summaryAmount, summaryInputsAmount)
	}
	return inputInfos, summaryInputsAmount, nil
}
