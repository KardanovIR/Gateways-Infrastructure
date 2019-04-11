package router_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/clientgrpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"testing"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Router/grpc"
	blClient "github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/blockchain"
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

// TestRouter_ListenerRequests checks add and remove task for waves listener
// it can't run without running nginx on port 50000 and waves listener
func TestRouter_ListenerRequests(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	// client to adapters and listeners send requests to port of nginx
	if err := clientgrpc.NewUniversalAdapterClient(context.Background(), ":50000"); err != nil {
		log.Fatal("Can't init grpc clients: ", err)
	}
	requestCtx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("blockchain-type",
		"waves"))
	addTaskReply, err := clientgrpc.GetUniversalClient().AddTask(requestCtx, &blClient.AddTaskRequest{
		ListenTo: &blClient.ListenObject{Type: "Address", Value: "123456"},
		TaskType: "2",
	})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if len(addTaskReply.TaskId) == 0 {
		t.FailNow()
	}
	_, err = clientgrpc.GetUniversalClient().RemoveTask(requestCtx, &blClient.RemoveTaskRequest{TaskId: addTaskReply.TaskId})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
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
