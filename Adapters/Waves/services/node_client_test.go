package services

import (
	"context"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
	"strconv"
	"testing"
)

func TestNodeClient(t *testing.T) {
	ctx := context.Background()
	// setup
	log, _ := logger.Init(false, logger.DEBUG)
	err := config.Load("./testdata/config_test.yml")
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	err = New(ctx, config.Cfg.Node.Host)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	block, err := GetNodeClient().GetLastBlockHeight(ctx)
	if err != nil {
		t.Fail()
	}
	bl, err := strconv.Atoi(block)
	if bl < 533148 {
		t.Fail()
	}
}
