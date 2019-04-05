package clientgrpc

import (
	"context"
	"sync"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"google.golang.org/grpc"
)

//

var (
	rpcClientInstance sync.Once
	grpcCl            pb.AdapterClient
)

// New create grpc client with connection to grpc server
func New(ctx context.Context, host string) error {
	log := logger.FromContext(ctx)
	var err error
	rpcClientInstance.Do(func() {
		conn, e := grpc.Dial(host, grpc.WithInsecure())
		if e != nil {
			err = e
			return
		}
		grpcCl = pb.NewAdapterClient(conn)
	})
	if err != nil {
		log.Errorf("error during initialise node client: %s", err)
	}
	return err
}

func GetClient() pb.AdapterClient {
	rpcClientInstance.Do(func() {
		panic("try to get node client before it's creation!")
	})
	return grpcCl
}
