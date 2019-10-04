package services

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
)

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

	amount := 0.0
	for _, output := range tx.Vout {
		if len(output.ScriptPubKey.Addresses) == 0 {
			continue
		}
		inputs = append(inputs, models.InputOutputInfo{
			Amount:  fmt.Sprintf("%f", output.Value),
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

func (cl *nodeClient) Fee(ctx context.Context) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'Fee'")
	// todo real fee calc or use parameter
	fee := MinOutputValue * 4
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
