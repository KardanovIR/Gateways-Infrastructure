package server

import (
	"context"
	"strconv"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
)

// get balance for address with assets balances
func (s *grpcServer) GetAllBalances(ctx context.Context, in *pb.AddressRequest) (*pb.GetAllBalancesReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetAllBalances for address %s", in.Address)
	var balance, err = s.nodeClient.GetAllBalances(ctx, in.Address)
	if err != nil {
		log.Errorf("getting all balances fails: %s", err)
		return nil, err
	}
	wb := strconv.FormatUint(balance.Amount, 10)
	assetBalances := make([]*pb.GetAllBalancesReply_AssetBalance, 0, len(balance.Assets))
	for c, amount := range balance.Assets {
		assetBalances = append(assetBalances,
			&pb.GetAllBalancesReply_AssetBalance{AssetId: c, Amount: strconv.FormatUint(amount, 10)},
		)
	}
	result := pb.GetAllBalancesReply{Amount: wb, AssetBalances: assetBalances}
	return &result, nil
}
