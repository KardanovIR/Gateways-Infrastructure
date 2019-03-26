package service

import (
	"context"
	"sync"

	"github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/ethAdapter"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/ethListener"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/wavesAdapter"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/wavesListener"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/model"
)

var (
	serviceSync     sync.Once
	serviceInstance IBlockChainsService
)

type IBlockChainsService interface {
	GetTransactionStatus(ctx context.Context, blockchain model.Blockchain, txID string) (string, error)
}

type blockchainsService struct {
	ethAdapter    ethAdapter.CommonClient
	wavesAdapter  wavesAdapter.CommonClient
	ethListener   ethListener.ListenerClient
	wavesListener wavesListener.ListenerClient
}

func New(ethAdapter ethAdapter.CommonClient,
	wavesAdapter wavesAdapter.CommonClient,
	ethListener ethListener.ListenerClient,
	wavesListener wavesListener.ListenerClient) IBlockChainsService {
	serviceSync.Do(func() {
		serviceInstance = &blockchainsService{
			ethAdapter:    ethAdapter,
			wavesAdapter:  wavesAdapter,
			ethListener:   ethListener,
			wavesListener: wavesListener,
		}
	})
	return serviceInstance
}
