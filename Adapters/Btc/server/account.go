package server

import (
	"context"
	"strconv"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
)

// get balance for address with assets balances
func (s *grpcServer) GetAllBalances(ctx context.Context, in *pb.AddressRequest) (*pb.GetAllBalanceReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetAllBalances for address %s", in.Address)
	var balance, err = s.nodeClient.GetAllBalances(ctx, in.Address)
	if err != nil {
		log.Errorf("getting all balances fails: %s", err)
		return nil, err
	}
	wb := strconv.FormatUint(balance.Amount, 10)
	assetBalances := make([]*pb.GetAllBalanceReply_AssetBalance, 0, len(balance.Assets))
	for c, amount := range balance.Assets {
		assetBalances = append(assetBalances,
			&pb.GetAllBalanceReply_AssetBalance{AssetId: c, Amount: strconv.FormatUint(amount, 10)},
		)
	}
	result := pb.GetAllBalanceReply{Amount: wb, AssetBalances: assetBalances}
	return &result, nil
}

func (s *grpcServer) GetAllBalance(ctx context.Context, in *pb.GetAllBalanceRequest) (*pb.GetAllBalanceReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetAllBalances for address %s", in.Address)
	var balance, err = s.nodeClient.GetAllBalances(ctx, in.Address)
	if err != nil {
		log.Errorf("getting all balances fails: %s", err)
		return nil, err
	}
	wb := strconv.FormatUint(balance.Amount, 10)
	assetBalances := make([]*pb.GetAllBalanceReply_AssetBalance, 0, len(balance.Assets))
	for c, amount := range balance.Assets {
		assetBalances = append(assetBalances,
			&pb.GetAllBalanceReply_AssetBalance{AssetId: c, Amount: strconv.FormatUint(amount, 10)},
		)
	}
	result := pb.GetAllBalanceReply{Amount: wb, AssetBalances: assetBalances}
	return &result, nil
}
