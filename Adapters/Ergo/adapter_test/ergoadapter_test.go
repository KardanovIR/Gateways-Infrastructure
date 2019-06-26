package adapter_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/config"
	adapter "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/services"
	"google.golang.org/grpc"
)

var grpcCl adapter.AdapterClient

func TestGrpcClient_TransactionByHash(t *testing.T) {
	ctx, _ := beforeTests()

	txInfo, err := grpcCl.TransactionByHash(ctx, &adapter.TransactionByHashRequest{TxHash: "1857c4e2490ff80cec9dc2ffdf64fb367744130c39641106562f88cf696f5096"})
	assert.Nil(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, 1, len(txInfo.Inputs))
	assert.Equal(t, 1, len(txInfo.Outputs))
	assert.Equal(t, "1200000", txInfo.Outputs[0].Amount)
	assert.Equal(t, "3WwgqLMBZhUWVHQUoYakSmcJwte8TPYM3gFkYeJ84S3NP21T2uJg", txInfo.Outputs[0].Address)

	assert.Equal(t, "2200000", txInfo.Inputs[0].Amount)
	assert.Equal(t, "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs", txInfo.Inputs[0].Address)

	assert.Equal(t, "1000000", txInfo.Fee)
	assert.Equal(t, "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs", txInfo.SenderAddress)
	assert.Equal(t, "3WwgqLMBZhUWVHQUoYakSmcJwte8TPYM3gFkYeJ84S3NP21T2uJg", txInfo.RecipientAddress)
	assert.Equal(t, "1857c4e2490ff80cec9dc2ffdf64fb367744130c39641106562f88cf696f5096", txInfo.TxHash)
	assert.Equal(t, "1200000", txInfo.Amount)
	assert.Equal(t, "SUCCESS", txInfo.Status)

	// unknown tx id
	txInfo2, err := grpcCl.TransactionByHash(ctx, &adapter.TransactionByHashRequest{TxHash: "1857c4e2490ff80cec9dc2ffdf64fb367744130c39641106562f88cf696f5000"})
	assert.Nil(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, "UNKNOWN", txInfo2.Status)

}

func beforeTests() (context.Context, logger.ILogger) {
	ctx := context.Background()
	log, _ := logger.Init(false, logger.DEBUG)
	err := config.Load("./testdata/config_adapter_test.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = services.New(ctx, config.Cfg.Node)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := server.InitAndStart(ctx, config.Cfg.Port, services.GetNodeClient()); err != nil {
			log.Fatal("Can't start grpc server", err)
		}
	}()

	conn, e := grpc.Dial(":"+config.Cfg.Port, grpc.WithInsecure())
	if e != nil {
		log.Fatal(e)
	}
	grpcCl = adapter.NewAdapterClient(conn)

	return ctx, log
}
