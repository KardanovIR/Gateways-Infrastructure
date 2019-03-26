package clientgrpc

import (
	"context"
	"google.golang.org/grpc"
	"sync"

	pbEthAdapter "github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/ethAdapter"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
)

var (
	EthAdapterSync   sync.Once
	ethAdapterClient pbEthAdapter.CommonClient
)

// New create grpc eth adapter client with connection to grpc server
func NewEthAdapterClient(ctx context.Context, host string) error {
	log := logger.FromContext(ctx)
	var err error
	EthAdapterSync.Do(func() {
		log.Infof("setup connection to eth adapter %s", host)
		conn, e := grpc.Dial(host, grpc.WithInsecure())
		if e != nil {
			err = e
			log.Errorf("setup connection to eth adapter fails: %s", err)
			return
		}
		ethAdapterClient = pbEthAdapter.NewCommonClient(conn)
	})
	return err
}

func GetEthAdapterClient() pbEthAdapter.CommonClient {
	EthAdapterSync.Do(func() {
		panic("try to get eth adapter client before it's creation!")
	})
	return ethAdapterClient
}
