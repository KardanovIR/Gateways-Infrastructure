package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/repositories"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/services/node_client"

	"testing"
)

const (
	privateKey1      = "cf6113f8ab2bd2e4cd9f89c40dcb4bcd5e07f635f87ba8bd53ff611a710c43dd"
	address1         = "mzrDT1HkUV6gBDa1rMDkXKy37wedV8N8ve"
	privateKey2      = "d47cbc5bbaa5bb4fc92c99f6064cb0fb9dd7aba78d51018cf0b438094ae8f1ea"
	address2         = "mwjKGKKxTaNNwnCAhGLvPXTR8Mn6P21aP1"
	recipientAddress = "2Mxd3wMiJEhHqcMPX8BrFwHxXSSsDvrrpJN"
)

func TestNodeClient_SendTx(t *testing.T) {
	ctx, log := beforeTest()
	outs := make([]*models.Output, 2)
	outs[0] = &models.Output{Address: recipientAddress, Amount: 80000}
	outs[1] = &models.Output{Address: "2MwWZAgprCccvQ7LmN4YUPPDd6Jn4FgKSe8", Amount: 210000}
	raw, err := node_client.GetNodeClient().CreateRawTx(ctx, []string{address1, address2}, address1, outs)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	privates := make(map[string]string)
	privates[address1] = privateKey1
	privates[address2] = privateKey2
	signedTx, err := node_client.GetNodeClient().SignTransaction(ctx, raw, privates)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}

	txId, err := node_client.GetNodeClient().SendTransaction(ctx, signedTx)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	log.Info("txId = ", txId)
	assert.True(t, len(txId) > 0)
}

func TestNodeClient_TransactionByHash(t *testing.T) {
	ctx, log := beforeTest()
	tx, err := node_client.GetNodeClient().TransactionByHash(ctx, "4d996fd7524c46a3afe959d0823ab15f30312bbf02e86a815749bca163b2664c")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, "2Mxd3wMiJEhHqcMPX8BrFwHxXSSsDvrrpJN", tx.From)
	assert.Equal(t, "61100", tx.Fee)
	assert.Equal(t, models.TxStatusSuccess, tx.Status)
	assert.Equal(t, "4d996fd7524c46a3afe959d0823ab15f30312bbf02e86a815749bca163b2664c", tx.TxHash)
	assert.Equal(t, "2071140", tx.Amount)
	assert.Equal(t, 2, len(tx.Outputs))
	assert.Equal(t, 1, len(tx.Inputs))
	assert.Equal(t, "2132240", tx.Inputs[0].Amount)
	assert.Equal(t, "2Mxd3wMiJEhHqcMPX8BrFwHxXSSsDvrrpJN", tx.Inputs[0].Address)
	assert.Equal(t, "1000000", tx.Outputs[0].Amount)
	assert.Equal(t, "mwjKGKKxTaNNwnCAhGLvPXTR8Mn6P21aP1", tx.Outputs[0].Address)
	assert.Equal(t, "1071140", tx.Outputs[1].Amount)
	assert.Equal(t, "2Mxd3wMiJEhHqcMPX8BrFwHxXSSsDvrrpJN", tx.Outputs[1].Address)
	txUnknown, err := node_client.GetNodeClient().TransactionByHash(ctx, "4d996fd7524c46a3afe959d0823ab15f30312bbf02e86a815749bca163b26641")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, models.TxStatusUnKnown, txUnknown.Status)
}

func beforeTest() (context.Context, logger.ILogger) {
	ctx := context.Background()
	log, _ := logger.Init(false, logger.DEBUG)
	err := config.Load("./testdata/config_test.yml")
	if err != nil {
		log.Fatal(err)
	}
	if err := repositories.New(ctx, config.Cfg.Db); err != nil {
		log.Fatal("can't create repository: ", err)
	}
	if err := node_client.New(ctx, config.Cfg.Node, repositories.GetRepository()); err != nil {
		log.Fatal("can't create node's client: ", err)
	}
	return ctx, log
}
