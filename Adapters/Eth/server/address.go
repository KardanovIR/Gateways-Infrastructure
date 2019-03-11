package server

import (
	"context"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
)

// Generate address
func (s *grpcServer) GenerateAddress(ctx context.Context, in *pb.EmptyRequest) (*pb.GenerateAddressReply, error) {
	log := logger.FromContext(ctx)
	log.Info("GenerateAddress")

	address, err := s.nodeClient.GenerateAddress(ctx)
	if err != nil {
		log.Errorf("address generation fails: %s", err)
		return nil, err
	}

	return &pb.GenerateAddressReply{Address: address}, nil
}

// Validate address
func (s *grpcServer) ValidateAddress(ctx context.Context, in *pb.AddressRequest) (*pb.ValidateAddressReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("ValidateAddress %s", in.Address)
	ok := s.nodeClient.IsAddressValid(ctx, in.Address)
	return &pb.ValidateAddressReply{Valid: ok}, nil
}
