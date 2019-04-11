package clientgrpc

import (
	"context"
	"google.golang.org/grpc"
	"sync"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/blockchain"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
)

var (
	universalSync   sync.Once
	universalClient UniversalGrpcClient
)

// New create grpc waves adapter client with connection to grpc server
func NewUniversalAdapterClient(ctx context.Context, host string) error {
	log := logger.FromContext(ctx)
	var err error
	universalSync.Do(func() {
		log.Infof("setup connection to proxy %s", host)
		conn, e := grpc.Dial(host, grpc.WithInsecure(), grpc.WithAuthority("service"))
		if e != nil {
			err = e
			log.Errorf("setup connection to proxy fails: %s", err)
			return
		}
		universalClient = UniversalGrpcClient{pb.NewAdapterClient(conn), pb.NewListenerClient(conn)}
	})
	return err
}

func GetUniversalClient() UniversalGrpcClient {
	universalSync.Do(func() {
		panic("try to get waves adapter client before it's creation!")
	})
	return universalClient
}

type UniversalGrpcClient struct {
	pb.AdapterClient
	pb.ListenerClient
}