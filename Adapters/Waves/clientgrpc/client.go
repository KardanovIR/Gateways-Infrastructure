package clientgrpc

import (
	"context"
	"google.golang.org/grpc"
	"sync"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/grpc"
)

var (
	rpcClientInstance sync.Once
	grpcCl            pb.AdapterClient
)

// New create grpc client with connection to grpc server
func New(ctx context.Context, host string) error {
	var err error
	rpcClientInstance.Do(func() {
		conn, e := grpc.Dial(host, grpc.WithInsecure())
		if e != nil {
			err = e
			return
		}
		grpcCl = pb.NewAdapterClient(conn)
	})
	return err
}

func GetClient() pb.AdapterClient {
	rpcClientInstance.Do(func() {
		panic("try to get node client before it's creation!")
	})
	return grpcCl
}
