package service

import (
	"context"
	"sync"

	"github.com/wavesplatform/GatewaysInfrastructure/Router/clientgrpc"
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
	universal clientgrpc.UniversalGrpcClient
}

func New(universal clientgrpc.UniversalGrpcClient) IBlockChainsService {
	serviceSync.Do(func() {
		serviceInstance = &blockchainsService{
			universal: universal,
		}
	})
	return serviceInstance
}
