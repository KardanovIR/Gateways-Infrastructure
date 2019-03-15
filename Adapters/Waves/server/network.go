package server

import (
	"context"
	"strconv"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
)

func (s *grpcServer) GetLastBlockHeight(ctx context.Context, in *pb.BlockRequest) (*pb.BlockReply, error) {
	log := logger.FromContext(ctx)
	log.Info("GetLastBlockHeight")

	var blockHeight, err = s.nodeClient.GetLastBlockHeight(ctx)
	if err != nil {
		log.Errorf("get last block height fails: %s", err)
		return nil, err
	}

	return &pb.BlockReply{Block: blockHeight}, nil
}

// Get suggested transaction's fee
func (s *grpcServer) Fee(ctx context.Context, in *pb.EmptyRequest) (*pb.FeeReply, error) {
	log := logger.FromContext(ctx)
	log.Info("Fee request")

	var fee, err = s.nodeClient.Fee(ctx)
	if err != nil {
		log.Errorf("get fee fails: %s", err)
		return nil, err
	}
	f, err := strconv.FormatUint(fee, 10), nil
	if err != nil {
		log.Errorf("convert fee fails: %s", err)
		return nil, err
	}
	return &pb.FeeReply{Fee: f}, nil
}
