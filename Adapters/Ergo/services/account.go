package services

import (
	"context"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/models"
)

func (cl *nodeClient) GetAllBalances(ctx context.Context, address string) (*models.AccountBalance, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'GetBalanceWithAssets' for address %s", address)
	balance := models.AccountBalance{}
	// todo implementation
	return &balance, nil
}
