package services

import (
	"context"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
)

func (cl *nodeClient) GetAllBalances(ctx context.Context, address string) (*models.Balance, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'GetAllBalances' for address %s", address)
	balances, err := cl.rep.GetBalanceForAddresses(ctx, []string{address})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, b := range balances {
		if b.Address == address {
			return &b, nil
		}
	}
	return &models.Balance{}, nil
}
