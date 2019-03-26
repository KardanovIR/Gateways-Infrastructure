package clientgrpc

import (
	"context"
	"google.golang.org/grpc"
	"sync"

	pbEthL "github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/ethListener"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
)

var (
	ethListenerSync   sync.Once
	ethListenerClient pbEthL.ListenerClient
)

// New create grpc eth listener client with connection to grpc server
func NewEthListenerClient(ctx context.Context, host string) error {
	log := logger.FromContext(ctx)
	var err error
	ethListenerSync.Do(func() {
		log.Infof("setup connection to eth listener %s", host)
		conn, e := grpc.Dial(host, grpc.WithInsecure())
		if e != nil {
			err = e
			log.Errorf("setup connection to eth listener fails: %s", err)
			return
		}
		ethListenerClient = pbEthL.NewListenerClient(conn)
	})
	return err
}

func GetEthListenerClient() pbEthL.ListenerClient {
	ethListenerSync.Do(func() {
		panic("try to get eth listener client before it's creation!")
	})
	return ethListenerClient
}
