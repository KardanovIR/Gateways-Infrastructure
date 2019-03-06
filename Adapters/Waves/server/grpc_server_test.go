package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"sync"
	"testing"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
)

const serverPort = "20001"

var oneServer sync.Once

func TestGrpcServerGetLastBlockHeight(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	conn, err := startServerAndGetConnection(ctx)
	if err != nil {
		log.Error("connection to grpc server don't establish", err)
		t.Fail()
	}
	defer conn.Close()

	c := pb.NewCommonClient(conn)
	reply, err := c.GetLastBlockHeight(ctx, &pb.BlockRequest{})
	if err != nil {
		log.Error("getting last block height fails", err)
		t.Fail()
	}
	if reply.Block != "3" {
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

func (n *nodeClientMock) GetLastBlockHeight(ctx context.Context) (string, error) {
	return "3", nil
}
