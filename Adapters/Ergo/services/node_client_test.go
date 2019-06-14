package services

import (
	"context"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
)

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
