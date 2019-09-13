package server

import (
	"context"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
)

// Validate address
func (s *grpcServer) ValidateAddress(ctx context.Context, in *pb.AddressRequest) (*pb.ValidateAddressReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("ValidateAddress %s", in.Address)
	var ok, err = s.nodeClient.ValidateAddress(ctx, in.Address)
	if err != nil {
		log.Debugf("validate address fails: %s", err)
	}
	return &pb.ValidateAddressReply{Valid: ok}, nil
}
