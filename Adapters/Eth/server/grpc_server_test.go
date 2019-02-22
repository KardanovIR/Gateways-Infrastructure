package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"math/big"
	"sync"
	"testing"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
)

const serverPort = "20001"

var oneServer sync.Once

func TestGrpcServerGasPrice(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	conn, err := startServerAndGetConnection(ctx)
	if err != nil {
		log.Error("connection to grpc server don't establish", err)
		t.Fail()
	}
	defer conn.Close()

	c := pb.NewCommonClient(conn)
	reply, err := c.GasPrice(ctx, &pb.GasPriceRequest{})
	if err != nil {
		log.Error("getting gas price fails", err)
		t.Fail()
	}
	if reply.GasPrice != "3" {
		t.Fail()
	}
}

func startServerAndGetConnection(ctx context.Context) (*grpc.ClientConn, error) {
	log := logger.FromContext(ctx)
	oneServer.Do(func() {
		go func() {
			err := InitAndStart(ctx, serverPort, &nodeClientMock{})
			if err != nil {
				log.Fatal("server can't start", err)
			}
		}()
	})
	return grpc.Dial(fmt.Sprint(":", serverPort), grpc.WithInsecure())
}

type nodeClientMock struct {
}

func (n *nodeClientMock) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return new(big.Int).SetInt64(3), nil
}
