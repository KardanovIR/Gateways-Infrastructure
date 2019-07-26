package server

import (
	"context"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/server/converter"
)

func (s *grpcServer) GasPrice(ctx context.Context, in *pb.GasPriceRequest) (*pb.GasPriceReply, error) {
	log := logger.FromContext(ctx)
	log.Info("GasPrice")

	var gasPrice, err = s.nodeClient.SuggestGasPrice(ctx)
	if err != nil {
		log.Errorf("getting gas price fails: %s", err)
		return nil, err
	}

	return &pb.GasPriceReply{GasPrice: converter.ToTargetAmountStr(gasPrice)}, nil
}

// Get suggested transaction's fee
func (s *grpcServer) SuggestFee(ctx context.Context, in *pb.EmptyRequest) (*pb.SuggestFeeReply, error) {
	log := logger.FromContext(ctx)
	log.Info("SuggestFee")

	var fee, err = s.nodeClient.SuggestFee(ctx)
	if err != nil {
		log.Errorf(" getting fee fails: %s", err)
		return nil, err
	}

	return &pb.SuggestFeeReply{Fee: converter.ToTargetAmountStr(fee)}, nil
}
