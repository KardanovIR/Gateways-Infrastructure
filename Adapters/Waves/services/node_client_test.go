package services

import (
	"context"
	"strconv"
	"testing"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
)

func TestNodeClient(t *testing.T) {
	ctx, _ := beforeTest()
	// setup

	block, err := GetNodeClient().GetLastBlockHeight(ctx)
	if err != nil {
		t.Fail()
	}
	bl, err := strconv.Atoi(block)
	if bl < 533148 {
		t.Fail()
	}
}

func beforeTest() (context.Context, logger.ILogger) {
	ctx := context.Background()
	log, _ := logger.Init(false, logger.DEBUG)
	err := config.Load("./testdata/config_test.yml")
	if err != nil {
		log.Fatal(err)
	}
	if err := New(ctx, config.Cfg.Node); err != nil {
		log.Fatal("can't create node's client: ", err)
	}
	return ctx, log
}
