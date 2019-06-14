package server

import (
	"context"
	"strconv"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
)

// Get transaction's fee
func (s *grpcServer) Fee(ctx context.Context, in *pb.FeeRequest) (*pb.FeeReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("Fee request for sender %s, assetId %s", in.SendersPublicKey, in.AssetId)
	var fee, err = s.nodeClient.Fee(ctx, in.SendersPublicKey, in.AssetId)
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
