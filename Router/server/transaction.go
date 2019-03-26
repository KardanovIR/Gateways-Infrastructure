package server

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Router/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/model"
)

func (s *grpcServer) GetTransactionStatus(ctx context.Context, in *pb.GetTransactionStatusRequest) (*pb.GetTransactionStatusReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetTransactionStatus %+v", in)
	if len(in.TxId) == 0 {
		err := errors.New("parameter 'txId' can't be empty")
		log.Error(err)
		return nil, err
	}
	b := model.Blockchain(in.Blockchain)
	if !b.Exist() {
		err := fmt.Errorf("don't support blockchain '%s'", b)
		log.Error(err)
		return nil, err
	}
	var status, err = s.service.GetTransactionStatus(ctx, b, in.TxId)
	if err != nil {
		log.Errorf("getting transaction's status fails: %s", err)
		return nil, err
	}
	return &pb.GetTransactionStatusReply{Status: status}, nil
}
