package router_test

import (
	"context"
	"google.golang.org/grpc"
	"testing"

	"github.com/magiconair/properties/assert"
	pb "github.com/wavesplatform/GatewaysInfrastructure/Router/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
)

// TestRouter_TransactionStatus checks getting status of waves and ethereum transactions
// it can't run without running router, waves and eth adapters
func TestRouter_TransactionStatus(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	cl, err := routerClient(":5010")
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

func routerClient(host string) (pb.RouterClient, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewRouterClient(conn), nil
}
