package adapter_test

import (
	"context"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/clientgrpc"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/config"
	adapter "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/services"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/services/converter"
)

var grpcCl adapter.AdapterClient

func TestGrpcClient_TransactionByHash(t *testing.T) {
	ctx, _ := beforeTests()

	txInfo, err := grpcCl.TransactionByHash(ctx, &adapter.TransactionByHashRequest{TxHash: "e3c732ac6d902dbe1b67825d7dba0ae16636708e61a0b94bf44c2550e86e7f62"})
	assert.Nil(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, 1, len(txInfo.Inputs))
	assert.Equal(t, 1, len(txInfo.Outputs))
	assert.Equal(t, "400000", txInfo.Outputs[0].Amount)
	assert.Equal(t, "9f1a5ZP3VWs6WGo32imb8aaPT6uNrzCqwD7wwTR7SYphfDWA1km", txInfo.Outputs[0].Address)

	assert.Equal(t, "500000", txInfo.Inputs[0].Amount)
	assert.Equal(t, "9eaFpf4DR1Fj3WnCvDdgfNNdfa8tAZ1Ga21YchCZpeFSEFtkKDq", txInfo.Inputs[0].Address)

	assert.Equal(t, "100000", txInfo.Fee)
	assert.Equal(t, "9eaFpf4DR1Fj3WnCvDdgfNNdfa8tAZ1Ga21YchCZpeFSEFtkKDq", txInfo.SenderAddress)
	assert.Equal(t, "9f1a5ZP3VWs6WGo32imb8aaPT6uNrzCqwD7wwTR7SYphfDWA1km", txInfo.RecipientAddress)
	assert.Equal(t, "e3c732ac6d902dbe1b67825d7dba0ae16636708e61a0b94bf44c2550e86e7f62", txInfo.TxHash)
	assert.Equal(t, "400000", txInfo.Amount)
	assert.Equal(t, "SUCCESS", txInfo.Status)

	// unknown tx id
	txInfo2, err := grpcCl.TransactionByHash(ctx, &adapter.TransactionByHashRequest{TxHash: "1857c4e2490ff80cec9dc2ffdf64fb367744130c39641106562f88cf696f5000"})
	assert.Nil(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, "UNKNOWN", txInfo2.Status)

}

func TestGrpcClient_AccountBalance(t *testing.T) {
	ctx, _ := beforeTests()

	balance, err := grpcCl.GetAllBalances(ctx, &adapter.AddressRequest{Address: "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs"})
	assert.Nil(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, "92430000", balance.Amount)
}

func beforeTests() (context.Context, logger.ILogger) {
	ctx := context.Background()
	log, _ := logger.Init(false, logger.DEBUG)
	err := config.Load("./testdata/config_adapter_test.yml")
	if err != nil {
		log.Fatal(err)
	}
	converter.Init(ctx, config.Cfg.Decimals)
	err = services.New(ctx, config.Cfg.Node)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := server.InitAndStart(ctx, config.Cfg.Port, services.GetNodeClient()); err != nil {
			log.Fatal("Can't start grpc server", err)
		}
	}()

	if err := clientgrpc.New(ctx, ":"+config.Cfg.Port); err != nil {
		log.Fatal("Can't init grpc client", err)
	}
	grpcCl = clientgrpc.GetClient()
	return ctx, log
}
