package server

import (
	"context"
	"errors"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
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

func (s *grpcServer) GetUnspentInputs(ctx context.Context, in *pb.AddressRequest) (*pb.GetUnspentInputsReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("ValidateAddress %s", in.Address)

	return nil, errors.New("not supported for this blockchain")
}