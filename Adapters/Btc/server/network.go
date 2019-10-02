package server

import (
	"context"
	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
)

// Get transaction's fee
func (s *grpcServer) Fee(ctx context.Context, in *pb.FeeRequest) (*pb.FeeReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("Fee request for sender %s, assetId %s", in.SendersPublicKey, in.AssetId)
	var fee, err = s.nodeClient.Fee(ctx, in.SendersPublicKey)
	if err != nil {
		log.Errorf("get fee fails: %s", err)
		return nil, err
	}
	return &pb.FeeReply{Fee: string(fee)}, nil
}