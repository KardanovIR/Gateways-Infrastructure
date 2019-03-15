package server

import (
	"context"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
)

// Generate address
func (s *grpcServer) GenerateAddress(ctx context.Context, in *pb.EmptyRequest) (*pb.GenerateAddressReply, error) {
	log := logger.FromContext(ctx)
	log.Info("GenerateAddress")
	var address, err = s.nodeClient.GenerateAddress(ctx)
	if err != nil {
		log.Errorf("generate address fails: %s", err)
		return nil, err
	}
	return &pb.GenerateAddressReply{Address: address}, nil
}

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
