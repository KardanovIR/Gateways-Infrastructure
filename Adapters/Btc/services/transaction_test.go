package services

import (
	"context"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/repositories"

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
	raw, err := GetNodeClient().CreateRawTx(ctx, []string{address1, address2}, address1, outs)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	privates := make(map[string]string)
	privates[address1] = privateKey1
	privates[address2] = privateKey2
	signedTx, err := GetNodeClient().SignTransaction(ctx, raw, privates)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}

	txId, err := GetNodeClient().SendTransaction(ctx, signedTx)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	log.Info("txId = ", txId)
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
	if err := New(ctx, config.Cfg.Node, repositories.GetRepository()); err != nil {
		log.Fatal("can't create node's client: ", err)
	}
	return ctx, log
}
