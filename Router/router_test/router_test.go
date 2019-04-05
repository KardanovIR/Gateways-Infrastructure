package router_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/clientgrpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/service"
	"google.golang.org/grpc"
	"log"
	"testing"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Router/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
)

// TestRouter_TransactionStatus checks getting status of waves and ethereum transactions
// it can't run without running router, waves and eth adapters
func TestRouter_TransactionStatus(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	cl, err := beforeTest()
	if err != nil {
		log.Fatal(err)
	}
	// tx in ropsten
	ethSt, err := cl.GetTransactionStatus(ctx, &pb.GetTransactionStatusRequest{
		TxId: "0x68392adbfd32cce6170eb909ad8c889319840593692df18c9a1b24818a1cfa1d", Blockchain: "ETH"})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, ethSt.Status, "SUCCESS")
	// waves tx in mainnet
	wavesSt, err := cl.GetTransactionStatus(ctx, &pb.GetTransactionStatusRequest{TxId: "699FRxjX1QpSpka7NWSUA7V5RdFiqf6F5nBANZxHmxep", Blockchain: "WAVES"})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, wavesSt.Status, "SUCCESS")

	wavesSt2, err := cl.GetTransactionStatus(ctx, &pb.GetTransactionStatusRequest{TxId: "699FRxjX1QpSpka7NWSUA7V5RdFiqf6F5nBANZx12xep", Blockchain: "WAVES"})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, wavesSt2.Status, "UNKNOWN")
}

func beforeTest() (pb.RouterClient, error) {
	ctx := context.Background()
	// client to adapters and listeners send requests to port of nginx
	if err := clientgrpc.NewUniversalAdapterClient(context.Background(), ":50000"); err != nil {
		log.Fatal("Can't init grpc clients: ", err)
	}
	blockchainService := service.New(clientgrpc.GetUniversalClient())

	// router's server starts on port 5555
	go func() {
		if err := server.InitAndStart(ctx, "5555", blockchainService); err != nil {
			log.Fatal("Can't start grpc server", err)
		}
	}()
	// client that send request to router
	conn, err := grpc.Dial(":5555", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return pb.NewRouterClient(conn), nil
}
