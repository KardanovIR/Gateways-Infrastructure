package server

import (
	"context"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	ok, msg, err := s.nodeClient.IsAddressValid(ctx, in.Address)
	if err != nil {
		return nil, err
	}
	if !ok && len(msg) > 0 {
		return &pb.ValidateAddressReply{Valid: ok}, status.Error(codes.InvalidArgument, msg)
	}
	return &pb.ValidateAddressReply{Valid: ok}, nil
}
