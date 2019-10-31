package server

import (
	"context"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
)

func (s *grpcServer) GasPrice(ctx context.Context, in *pb.GasPriceRequest) (*pb.GasPriceReply, error) {
	log := logger.FromContext(ctx)
	log.Info("GasPrice")

	var gasPrice, err = s.nodeClient.SuggestGasPrice(ctx)
	if err != nil {
		log.Errorf("getting gas price fails: %s", err)
		return nil, err
	}
	gas, err := s.converter.ToTargetAmountStr(ctx, gasPrice, "")
	if err != nil {
		log.Errorf("converting gas price fails: %s", err)
		return nil, err
	}
	return &pb.GasPriceReply{GasPrice: gas}, nil
}

// Get suggested transaction's fee
func (s *grpcServer) Fee(ctx context.Context, in *pb.FeeRequest) (*pb.SuggestFeeReply, error) {
	log := logger.FromContext(ctx)
	log.Info("SuggestFee")

	var fee, err = s.nodeClient.SuggestFee(ctx)
	if err != nil {
		log.Errorf(" getting fee fails: %s", err)
		return nil, err
	}
	return &pb.SuggestFeeReply{Fee: s.converter.ToCommissionStr(fee)}, nil
}
