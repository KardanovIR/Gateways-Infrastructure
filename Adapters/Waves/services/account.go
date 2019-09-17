package services

import (
	"context"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/models"
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

func (cl *nodeClient) GetAllBalances(ctx context.Context, address string) (*models.AccountBalance, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'GetBalanceWithAssets' for address %s", address)
	balance := models.AccountBalance{}
	wavesBalance, err := cl.GetBalance(ctx, address)
	if err != nil {
		log.Error("get waves balance fails", err)
		return nil, err
	}
	balance.Amount = wavesBalance
	a, err := proto.NewAddressFromString(address)
	if err != nil {
		log.Error("can't get address from string", err)
		return nil, err
	}
	assetBalances, _, err := cl.nodeClient.Assets.BalanceByAddress(ctx, a)
	if err != nil {
		log.Error("can't get assets balance", err)
		return nil, err
	}
	balance.Assets = make(map[string]uint64)
	for _, b := range assetBalances.Balances {
		balance.Assets[b.AssetId.String()] = b.Balance
	}
	return &balance, nil
}

// Feature implementation with asset filtering
//func (cl *nodeClient) GetAllBalances(ctx context.Context, address string, contracts ...string) (*models.AccountBalance, error) {
//	log := logger.FromContext(ctx)
//	log.Debugf("Call method 'GetAllBalances' for %s", address)
//	balances := models.AccountBalance{}
//	wavesBalance, err := cl.GetBalance(ctx, address)
//	if err != nil {
//		return nil, err
//	}
//	balances.Amount = wavesBalance
//	if len(contracts) == 0 {
//		return &balances, nil
//	}
//	a, err := proto.NewAddressFromString(address)
//	if err != nil {
//		log.Error("can't get address from string", err)
//		return nil, err
//	}
//	assetBalances, _, err := cl.nodeClient.Assets.BalanceByAddress(ctx, a)
//	if err != nil {
//		log.Error("can't get assets balance", err)
//		return nil, err
//	}
//	balances.Assets = make(map[string]uint64)
//	for _,c := range contracts {
//		for _,b := range assetBalances.Balances {
//			if(c == b.AssetId.String()) {
//				balances.Assets[b.AssetId.String()] = b.Balance
//			}
//		}
//	}
//	return &balances, nil
//}
