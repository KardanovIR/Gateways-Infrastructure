package node

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
)

func TestNodeClient_CreateInputs(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	nc := nodeClient{rep: &RepMock{file: "unspent_tx_1.json"}, conf: config.Node{ChainParams: &chaincfg.TestNet3Params}}
	out1 := models.Output{Address: "mzrDT1HkUV6gBDa1rMDkXKy37wedV8N8ve", Amount: 10000}
	out2 := models.Output{Address: "2Mxd3wMiJEhHqcMPX8BrFwHxXSSsDvrrpJN", Amount: 15000}
	outs := []*models.Output{&out1, &out2}
	inputsInfos, changeAmount, err := nc.CreateInputs(ctx, 12000, outs)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, 1, len(inputsInfos))
	assert.Equal(t, "txHash3", inputsInfos[0].Input.Txid)
	assert.Equal(t, uint64(1892), changeAmount)

	// not enough funds
	outs2 := []*models.Output{
		{Address: "mzrDT1HkUV6gBDa1rMDkXKy37wedV8N8ve", Amount: 50000},
		{Address: "2Mxd3wMiJEhHqcMPX8BrFwHxXSSsDvrrpJN", Amount: 49000},
	}
	_, _, err = nc.CreateInputs(ctx, 12000, outs2)
	assert.NotNil(t, err)

	outs3 := []*models.Output{
		{Address: "mzrDT1HkUV6gBDa1rMDkXKy37wedV8N8ve", Amount: 10000},
		{Address: "2Mxd3wMiJEhHqcMPX8BrFwHxXSSsDvrrpJN", Amount: 48000},
	}
	inputsInfos3, changeAmount3, err := nc.CreateInputs(ctx, 12000, outs3)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}

	assert.Equal(t, 2, len(inputsInfos3))
	assert.Equal(t, "txHash4", inputsInfos3[0].Input.Txid)
	assert.Equal(t, "txHash3", inputsInfos3[1].Input.Txid)
	assert.Equal(t, uint64(7104), changeAmount3)
}

type RepMock struct {
	file string
}

func (r *RepMock) GetUnspentTxListForAddresses(ctx context.Context, addresses []string) ([]models.UnspentTx, error) {
	txInputsList := make([]models.UnspentTx, 0)
	inputsFomFile, err := ioutil.ReadFile("./testdata/" + r.file)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(inputsFomFile, &txInputsList); err != nil {
		return nil, err
	}
	return txInputsList, nil
}

func (*RepMock) GetUnspentTxByTxHashAndOutputNumber(ctx context.Context, txHash string, outputN uint32) (*models.UnspentTx, error) {
	panic("implement me")
}

func (*RepMock) GetBalanceForAddresses(ctx context.Context, addresses []string) ([]models.Balance, error) {
	panic("implement me")
}
