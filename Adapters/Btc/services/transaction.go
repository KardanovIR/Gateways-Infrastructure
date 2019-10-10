package services

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/converter"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
)

const RPC_INVALID_ADDRESS_OR_KEY = "-5"

type SendTxResponse struct {
	ID string `json:"id"`
}

func (cl *nodeClient) SignTransaction(ctx context.Context, txUnsigned []byte, privateKeyForAddress map[string]string) (txSigned []byte, err error) {
	log := logger.FromContext(ctx)
	log.Info("sign transaction")
	txRawInfo := TxRawInfo{}
	if err := json.Unmarshal(txUnsigned, &txRawInfo); err != nil {
		log.Error(err)
		return nil, err
	}
	tx, err := Deserialize(txRawInfo.RawTx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for inputIndex, in := range tx.TxIn {
		// find address for input
		inputAddr := ""
		for _, info := range txRawInfo.InputsInfo {
			if info.Input.Txid == in.PreviousOutPoint.Hash.String() && info.Input.Vout == in.PreviousOutPoint.Index {
				inputAddr = info.Address
				break
			}
		}
		if len(inputAddr) < 0 {
			err := fmt.Errorf("don't find address for input %s N = %d", in.PreviousOutPoint.Hash.String(), in.PreviousOutPoint.Index)
			log.Error(err)
			return nil, err
		}

		address, err := btcutil.DecodeAddress(inputAddr, cl.conf.ChainParams)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		previousOutputScript, err := txscript.PayToAddrScript(address)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		pk := privateKeyForAddress[inputAddr]
		prBytes, _ := hex.DecodeString(pk)
		privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), prBytes)
		sigScript, err := txscript.SignatureScript(tx, inputIndex, previousOutputScript, txscript.SigHashAll, privateKey, true)
		if err != nil {
			return nil, err
		}
		in.SignatureScript = sigScript
	}

	txSigned, err = Serialize(tx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return txSigned, nil
}

func (cl *nodeClient) SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'SendRawTransaction'")
	tx, err := Deserialize(txSigned)
	if err != nil {
		log.Error(err)
		return "", err
	}
	hash, err := cl.nodeClient.SendRawTransaction(tx, true)
	if err != nil {
		log.Errorf("sending transaction to node fails %s", err)
		return "", err
	}
	return hash.String(), nil
}

func (cl *nodeClient) TransactionByHash(ctx context.Context, txId string) (*models.TxInfo, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'TransactionByHash' for txID %s", txId)
	txHash, err := chainhash.NewHashFromStr(txId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	nodeTx, err := cl.nodeClient.GetRawTransactionVerbose(txHash)
	if err != nil {
		if strings.Contains(err.Error(), RPC_INVALID_ADDRESS_OR_KEY) {
			return &models.TxInfo{Status: models.TxStatusUnKnown}, nil
		}
		log.Error(err)
		return nil, err
	}
	return cl.parseTx(ctx, nodeTx)
}

func (cl *nodeClient) parseTx(ctx context.Context, tx *btcjson.TxRawResult) (*models.TxInfo, error) {
	log := logger.FromContext(ctx)
	inputs := make([]models.InputOutput, 0)
	outputs := make([]models.InputOutput, 0)
	inputAmountSum := uint64(0)
	inputsTxMap := make(map[string]*btcjson.TxRawResult) // keep tx which get from node to avoid double requests for one tx
	for _, input := range tx.Vin {
		// get tx which was used for input
		var inputTx *btcjson.TxRawResult
		if tx, ok := inputsTxMap[input.Txid]; ok {
			inputTx = tx
		} else {
			inputHash, err := chainhash.NewHashFromStr(input.Txid)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			tx, err := cl.nodeClient.GetRawTransactionVerbose(inputHash)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			inputTx = tx
			inputsTxMap[input.Txid] = tx
		}
		// find output which was used for input of current tx
		var vOut btcjson.Vout
		for _, previousOut := range inputTx.Vout {
			if previousOut.N == input.Vout {
				vOut = previousOut
				break
			}
		}
		if len(vOut.ScriptPubKey.Hex) == 0 {
			// not found output for input
			return nil, fmt.Errorf("not found output for input %s N = %d", input.Txid, input.Vout)
		}
		amount, err := converter.GetIntFromFloat(vOut.Value)
		if err != nil {
			return nil, err
		}
		inputAmountSum += amount
		address := ""
		if len(vOut.ScriptPubKey.Addresses) > 0 {
			address = vOut.ScriptPubKey.Addresses[0]
		}
		inputs = append(inputs, models.InputOutput{
			Value:   amount,
			Address: address,
		})
	}

	outputAmountSum := uint64(0)
	for _, output := range tx.Vout {
		if len(output.ScriptPubKey.Addresses) != 1 {
			// addresses count > 1 can be for multisign account (we haven't them)
			continue
		}
		amount, err := converter.GetIntFromFloat(output.Value)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, models.InputOutput{
			Value:   amount,
			Address: output.ScriptPubKey.Addresses[0],
		})
		outputAmountSum += amount
	}
	// sum different inputs(outputs) with the same address
	inputs = summarizeAmountByAddress(inputs)
	outputs = summarizeAmountByAddress(outputs)
	fee := inputAmountSum - outputAmountSum
	txInputs := make([]models.InputOutputInfo, len(inputs))
	txOutputs := make([]models.InputOutputInfo, len(outputs))
	for i, in := range inputs {
		txInputs[i] = models.InputOutputInfo{
			Amount:  strconv.FormatUint(in.Value, 10),
			Address: in.Address,
		}
	}
	for i, out := range outputs {
		txOutputs[i] = models.InputOutputInfo{
			Amount:  strconv.FormatUint(out.Value, 10),
			Address: out.Address,
		}
	}
	sender := ""
	if len(txInputs) == 1 {
		sender = txInputs[0].Address
	}
	var status models.TxStatus
	if tx.Confirmations == 0 {
		status = models.TxStatusPending
	} else {
		status = models.TxStatusSuccess
	}
	return &models.TxInfo{
		From:    sender,
		Amount:  strconv.FormatUint(outputAmountSum, 10),
		TxHash:  tx.Txid,
		Fee:     strconv.FormatUint(fee, 10),
		Status:  status,
		Inputs:  txInputs,
		Outputs: txOutputs,
	}, nil
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
		result = append(result, models.InputOutput{Address: a.Address, Value: amount})
	}
	return result
}

func hasAddress(address string, list []models.InputOutput) bool {
	for _, t := range list {
		if address == t.Address {
			return true
		}
	}
	return false
}

func (cl *nodeClient) Fee(ctx context.Context) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'Fee'")
	// todo real fee calc or use parameter
	fee := MinBtcOutputValueConst * 2
	return fee, nil
}

func Serialize(tx *wire.MsgTx) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Deserialize(rawTx []byte) (*wire.MsgTx, error) {
	msgTx := new(wire.MsgTx)
	err := msgTx.Deserialize(bytes.NewReader(rawTx))
	if err != nil {
		return nil, err
	}
	return msgTx, nil
}
