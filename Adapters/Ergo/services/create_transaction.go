package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/services/converter"
)

const (
	getCurrentBlockUrl               = "/blocks?limit=1"
	getUnspentTxByAddressUrlTemplate = "/transactions/boxes/byAddress/unspent/%s"

	minerErgoTree = "1005040004000e36100204a00b08cd0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798ea02d192a39a8cc7a701730073011001020402d19683030193a38cc7b2a57300000193c2b2a57301007473027303830108cdeeac93b1a57304"
	txFee         = uint64(1000000)
)

var (
	ergoTreePrefix = []byte{0x00, 0x08, 0xcd}
)

type emptyObject struct{}

// CreateRawTx creates transaction
func (cl *nodeClient) CreateRawTx(ctx context.Context, addressFrom string, outs []*models.Output) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'CreateRawTx' from %s to %v",
		addressFrom, outs)
	if len(addressFrom) == 0 || len(outs) == 0 {
		return nil, fmt.Errorf("wrong parameters addressFrom %s or outs", addressFrom)
	}
	amount := uint64(0)
	for _, r := range outs {
		r.Amount = converter.ToNodeAmount(r.Amount)
		amount += r.Amount
	}
	// get unspent input from explorer
	unspentTxList, err := cl.requestUnSpentInputs(ctx, addressFrom)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	fee := txFee
	txInputsList, inputsAmount, err := createTxInputs(ctx, unspentTxList, amount+fee)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	outputs := make([]*models.Output, 0)
	// output for recipient
	outputs = append(outputs, outs...)
	if inputsAmount > amount+fee {
		// output for charge (senders address)
		outputs = append(outputs, &models.Output{Address: addressFrom, Amount: inputsAmount - amount - fee})
	}

	// get current height
	height, err := cl.getCurrentHeight(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	outputTxList := cl.createTxOutputs(ctx, outputs, fee, height)
	tx := models.UnSignedTx{
		Inputs:     txInputsList,
		DataInputs: make([]interface{}, 0),
		Outputs:    outputTxList,
	}

	txBinary, err := json.Marshal(tx)
	if err != nil {
		log.Error(err)
		return []byte{}, err
	}
	log.Debugf("created tx %+v", string(txBinary))
	return txBinary, nil
}

func createTxInputs(ctx context.Context, unspentTxList []models.UnSpentTx, neededAmount uint64) (
	txInputs []models.TxInput, summaryAmount uint64, err error) {

	log := logger.FromContext(ctx)
	log.Debugf("create inputs for amount %d", neededAmount)
	// sort by amount (begins from the least)
	sort.Slice(unspentTxList, func(i, j int) bool {
		return unspentTxList[i].Value < unspentTxList[j].Value
	})
	type boxInput struct {
		id     string
		amount uint64
	}
	amountsSum := uint64(0)
	boxForInputs := make([]*boxInput, 0)
	hasFunds := false
	for _, box := range unspentTxList {
		amountsSum += box.Value
		boxForInputs = append(boxForInputs, &boxInput{id: box.ID, amount: box.Value})
		if amountsSum >= neededAmount {
			hasFunds = true
			break
		}
	}
	if !hasFunds {
		return nil, 0, fmt.Errorf("insufficient funds: need %d, has %d", neededAmount, amountsSum)
	}
	txInputsList := make([]models.TxInput, 0)
	inputsAmount := uint64(0)
	for _, b := range boxForInputs {
		txInputsList = append(txInputsList, models.TxInput{
			BoxId:         b.id,
			SpendingProof: models.SpendingProof{Extension: emptyObject{}},
		})
		inputsAmount += b.amount
	}
	return txInputsList, inputsAmount, nil
}

func (cl *nodeClient) createTxOutputs(ctx context.Context, outputs []*models.Output, fee, height uint64) []models.TxOutput {
	outputTxList := make([]models.TxOutput, 0)
	for _, o := range outputs {
		outputTxList = append(outputTxList, models.TxOutput{
			ErgoTree:            cl.ergoTreeByAddress(ctx, o.Address),
			Value:               o.Amount,
			CreationHeight:      height,
			Assets:              make([]interface{}, 0),
			AdditionalRegisters: emptyObject{},
		})
	}
	if fee > 0 {
		outputTxList = append(outputTxList, models.TxOutput{
			ErgoTree:            minerErgoTree,
			Value:               fee,
			CreationHeight:      height,
			Assets:              make([]interface{}, 0),
			AdditionalRegisters: emptyObject{},
		})
	}
	return outputTxList
}

func (cl *nodeClient) requestUnSpentInputs(ctx context.Context, address string) ([]models.UnSpentTx, error) {
	log := logger.FromContext(ctx)
	log.Infof("request unspent inputs for address %s", address)
	// /transactions/boxes/byAddress/unspent/${address}
	respUnspent, err := cl.Request(ctx, http.MethodGet,
		cl.conf.ExplorerUrl+fmt.Sprintf(getUnspentTxByAddressUrlTemplate, address), nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	unspentTxList := make([]models.UnSpentTx, 0)
	if err := json.Unmarshal(respUnspent, &unspentTxList); err != nil {
		log.Errorf("failed to get unspent inputs for address %s: %s", address, err)
		return nil, err
	}
	return unspentTxList, nil
}

func (cl *nodeClient) ergoTreeByAddress(ctx context.Context, address string) string {
	publicKey := cl.PublicKeyFromAddress(ctx, address)
	ergoTreePrefixLength := len(ergoTreePrefix)
	var ergoTreeBytes = make([]byte, len(publicKey)+ergoTreePrefixLength)
	// [0x00 0x08 0xcd public_key_bytes]
	copy(ergoTreeBytes[:ergoTreePrefixLength], ergoTreePrefix)
	copy(ergoTreeBytes[ergoTreePrefixLength:], publicKey)
	return hex.EncodeToString(ergoTreeBytes)
}
