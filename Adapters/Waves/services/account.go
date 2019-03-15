package services

import (
	"context"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func (cl *nodeClient) GetBalance(ctx context.Context, address string) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'GetBalance' for address %s", address)
	a, err := proto.NewAddressFromString(address)
	if err != nil {
		log.Error("can't get address from string", err)
		return 0, err
	}
	balance, _, err := cl.nodeClient.Addresses.Balance(ctx, a)
	if err != nil {
		log.Error("get balance fails", err)
		return 0, err
	}
	return balance.Balance, nil
}
