package clientgrpc

import (
	"context"
	"google.golang.org/grpc"
	"sync"

	pbWavesL "github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/wavesListener"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
)

var (
	wavesListSync       sync.Once
	wavesListenerClient pbWavesL.ListenerClient
)

// New create grpc waves listener client with connection to grpc server
func NewWavesListenerClient(ctx context.Context, host string) error {
	var err error
	log := logger.FromContext(ctx)
	wavesListSync.Do(func() {
		log.Infof("setup connection to waves listener %s", host)
		conn, e := grpc.Dial(host, grpc.WithInsecure())
		if e != nil {
			err = e
			log.Errorf("setup connection to waves listener fails: %s", err)
			return
		}
		wavesListenerClient = pbWavesL.NewListenerClient(conn)
	})
	return err
}

func GetWavesListenerClient() pbWavesL.ListenerClient {
	wavesListSync.Do(func() {
		panic("try to get waves listener client before it's creation!")
	})
	return wavesListenerClient
}
