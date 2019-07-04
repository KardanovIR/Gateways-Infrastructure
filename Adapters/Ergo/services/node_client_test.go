package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestGetNodeClient_ErgoTreeByAddress(t *testing.T) {
	// 0008cd0259a7f4ba52065bbb8ff7f5faa5b0bb372c2ab9008c633be0a7fe72aadabef6cd
	nc := nodeClient{}
	ergoTree := nc.ergoTreeByAddress(context.Background(), "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs")
	assert.Equal(t, "0008cd0259a7f4ba52065bbb8ff7f5faa5b0bb372c2ab9008c633be0a7fe72aadabef6cd", ergoTree)
}
