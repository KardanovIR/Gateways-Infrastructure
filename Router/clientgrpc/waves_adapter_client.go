package clientgrpc

import (
	"context"
	"google.golang.org/grpc"
	"sync"

	pbWavesA "github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/wavesAdapter"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
)

var (
	wavesAdapterSync   sync.Once
	wavesAdapterClient pbWavesA.CommonClient
)

// New create grpc waves adapter client with connection to grpc server
func NewWavesAdapterClient(ctx context.Context, host string) error {
	log := logger.FromContext(ctx)
	var err error
	wavesAdapterSync.Do(func() {
		log.Infof("setup connection to waves adapter %s", host)
		conn, e := grpc.Dial(host, grpc.WithInsecure())
		if e != nil {
			err = e
			log.Errorf("setup connection to waves adapter fails: %s", err)
			return
		}
		wavesAdapterClient = pbWavesA.NewCommonClient(conn)
	})
	return err
}

func GetWavesAdapterClient() pbWavesA.CommonClient {
	wavesAdapterSync.Do(func() {
		panic("try to get waves adapter client before it's creation!")
	})
	return wavesAdapterClient
}
