package node

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
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
	if len(changeAddress) == 0 || len(outs) == 0 {
		return nil, fmt.Errorf("wrong parameters changeAddress %s or outs", changeAddress)
	}
	amount := uint64(0)
	for _, o := range outs {
		if o.Amount < MinBtcOutputValueConst {
			return nil, fmt.Errorf("amount %d to sent to %s is less than min amount of fee", o.Amount, o.Address)
		}
		amount += o.Amount
	}
	feeRate, err := cl.FeeRateForKByte(ctx)
	if err != nil {
		log.Errorf("get fee rate fails %s", err)
		return nil, err
	}
	inputInfos, changeAmount, err := cl.CreateInputs(ctx, feeRate, outs)
	if err != nil {
		log.Errorf("%s", err)
		return nil, err
	}
	inputsForTx := make([]btcjson.TransactionInput, 0)
	for _, inputInfo := range inputInfos {
		inputsForTx = append(inputsForTx, inputInfo.Input)
	}
	destinations := make(map[btcutil.Address]btcutil.Amount)
	if changeAmount > MinBtcOutputValueConst {
		sendersAddress, err := btcutil.DecodeAddress(changeAddress, cl.conf.ChainParams)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		destinations[sendersAddress] = btcutil.Amount(changeAmount)
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
	// rawTx.TxOut
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

func (cl *nodeClient) CreateInputs(ctx context.Context, feeRate uint64, outs []*models.Output) (
	inputInfos []TxInputInfo, changeAmount uint64, err error) {
	log := logger.FromContext(ctx)
	txOutputs, err := cl.CreateOutputs(outs)
	if err != nil {
		log.Errorf("creation of outputs fails %s", err)
		return nil, 0, err
	}
	amount := uint64(0)
	for _, o := range outs {
		amount += o.Amount
	}
	// set empty addresses list ot get all inputs
	unspentInputs, err := cl.rep.GetUnspentTxListForAddresses(ctx, []string{})
	if err != nil {
		log.Errorf("get unspent inputs fails %s", err)
		return nil, 0, err
	}
	fee := uint64(0)
	var foundUnspentTxList []models.UnspentTx
	summaryInputsAmount := uint64(0)
	haveChangeOutput := false
	inputsCount := 1 // optimistic variant
	txSize := 0
	for {
		txSize = txsizes.EstimateSerializeSize(inputsCount, txOutputs, haveChangeOutput)
		newFee := cl.Fee(ctx, feeRate, txSize)
		log.Infof("estimate tx size %d; calculated fee for tx %d, fee on previous step %d", txSize, newFee, fee)
		if newFee <= fee {
			break
		}
		fee = newFee
		// if fee for new tx size is more than previous fee -> summary amount is changed -> recalculate inputs for new amount
		summaryAmount := amount + fee
		foundUnspentTxList = cl.findInputs(unspentInputs, summaryAmount)
		summaryInputsAmount = uint64(0)
		for _, input := range foundUnspentTxList {
			summaryInputsAmount += input.Amount
		}
		if summaryInputsAmount < summaryAmount {
			err = fmt.Errorf("insufficient funds: need %d, has %d", summaryAmount, summaryInputsAmount)
			log.Error(err)
			return nil, 0, err
		}
		changeAmount = summaryInputsAmount - (amount + fee)
		haveChangeOutput = changeAmount > MinBtcOutputValueConst
		inputsCount = len(foundUnspentTxList)
		log.Infof("inputs count for amount %d and fee %d = %d ", amount, fee, inputsCount)
	}
	log.Infof("result: tx inputs count %d; outputs count: %d;\n "+
		"has output for change: %v, change amount: %d; summary amount for recipients: %d; fee %d; inputs amount: %d; \n"+
		"calculated tx size (bytes) %d",
		len(foundUnspentTxList), len(txOutputs), haveChangeOutput, changeAmount, amount, fee, summaryInputsAmount, txSize)

	inputInfos = make([]TxInputInfo, 0)
	for _, input := range foundUnspentTxList {
		txInput := btcjson.TransactionInput{Txid: input.TxHash, Vout: input.OutputN}
		inputInfos = append(inputInfos, TxInputInfo{Address: input.Address, Input: txInput})
	}
	return inputInfos, changeAmount, nil
}

func (cl *nodeClient) CreateOutputs(outs []*models.Output) ([]*wire.TxOut, error) {
	wireOutList := make([]*wire.TxOut, 0)
	for _, out := range outs {
		outAddress, err := btcutil.DecodeAddress(out.Address, cl.conf.ChainParams)
		if err != nil {
			return nil, err
		}
		sourcePkScript, _ := txscript.PayToAddrScript(outAddress)
		wireOutList = append(wireOutList, wire.NewTxOut(int64(out.Amount), sourcePkScript))
	}
	return wireOutList, nil
}

func (cl *nodeClient) findInputs(unspentInputs []models.UnspentTx, targetAmount uint64) []models.UnspentTx {
	resultInputs := make([]models.UnspentTx, 0)
	// sort by amount (begins from the least)
	sort.Slice(unspentInputs, func(i, j int) bool {
		return unspentInputs[i].Amount < unspentInputs[j].Amount
	})
	for len(unspentInputs) > 0 {
		if targetAmount < MinBtcOutputValueConst {
			// rest of money - add one input (with the smallest amount, it amount can't be less than 1000)
			resultInputs = append(resultInputs, unspentInputs[0])
			return resultInputs
		}
		for _, input := range unspentInputs {
			if input.Amount >= targetAmount {
				// check difference between amounts (it will be change)
				// change (as every output) can't be less than 1000 - to avoid extra fee -> take next element with large amount
				if input.Amount == targetAmount || input.Amount-targetAmount > MinBtcOutputValueConst {
					// found suitable input
					resultInputs = append(resultInputs, input)
					return resultInputs
				}
			}
		}
		// didn't find suitable  - take last input with the largest amount
		lastIndex := len(unspentInputs) - 1
		input := unspentInputs[lastIndex]
		resultInputs = append(resultInputs, input)
		if input.Amount < targetAmount {
			targetAmount = targetAmount - input.Amount
		} else {
			targetAmount = input.Amount - targetAmount
		}
		// array without last element
		unspentInputs = unspentInputs[:lastIndex]
	}
	return resultInputs
}
