package server

import (
	"context"
	"strconv"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
)

// get balance for address
func (s *grpcServer) GetBalance(ctx context.Context, in *pb.AddressRequest) (*pb.GetBalanceReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetBalance for address %s", in.Address)
	var balance, err = s.nodeClient.GetBalance(ctx, in.Address)
	if err != nil {
		log.Errorf("getting balance fails: %s", err)
		return nil, err
	}
	b, err := strconv.FormatUint(balance, 10), nil
	if err != nil {
		log.Errorf("convert balance fails: %s", err)
		return nil, err
	}
	return &pb.GetBalanceReply{Balance: b}, nil
}
