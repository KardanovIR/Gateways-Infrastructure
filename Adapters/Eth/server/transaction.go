package server

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
)

func (s *grpcServer) GetRawTransaction(ctx context.Context, in *pb.RawTransactionRequest) (*pb.RawTransactionReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("getRawTransaction %+v", in)
	amount, ok := new(big.Int).SetString(in.Amount, 10)
	if !ok {
		err := fmt.Errorf("wrong amount value: %s", in.Amount)
		log.Error(err)
		return nil, err
	}
	var tx, err = s.nodeClient.CreateRawTransaction(ctx, in.AddressFrom, in.AddressTo, amount)
	if err != nil {
		log.Errorf("transaction's creation fails: %s", err)
		return nil, err
	}
	return &pb.RawTransactionReply{Tx: tx}, nil
}

// Sing transaction
func (s *grpcServer) SignTransaction(ctx context.Context, in *pb.SignTransactionRequest) (*pb.SignTransactionReply, error) {
	log := logger.FromContext(ctx)
	log.Info("SignTransaction")
	tx, err := s.nodeClient.SignTransaction(ctx, in.SenderAddress, []byte(in.Tx))
	if err != nil {
		log.Errorf("sign transaction fails: %s", err)
		return nil, err
	}
	return &pb.SignTransactionReply{Tx: tx}, nil
}

// Sing transaction by private key in parameters
func (s *grpcServer) SignTransactionWithPrivateKey(ctx context.Context, in *pb.SignTransactionWithPrivateKeyRequest) (*pb.SignTransactionReply, error) {
	log := logger.FromContext(ctx)
	log.Info("SignTransactionWithPrivateKey")
	if len(in.PrivateKey) == 0 {
		return nil, errors.New("private key can't be empty")
	}
	tx, err := s.nodeClient.SignTransactionWithPrivateKey(ctx, in.PrivateKey, []byte(in.Tx))
	if err != nil {
		log.Errorf("sign transaction fails: %s", err)
		return nil, err
	}
	return &pb.SignTransactionReply{Tx: tx}, nil
}

// Send transaction
func (s *grpcServer) SendTransaction(ctx context.Context, in *pb.SendTransactionRequest) (*pb.SendTransactionReply, error) {
	log := logger.FromContext(ctx)
	log.Info("SendTransaction")
	txHash, err := s.nodeClient.SendTransaction(ctx, []byte(in.Tx))
	if err != nil {
		log.Errorf("send transaction fails: %s", err)
		return nil, err
	}
	return &pb.SendTransactionReply{TxHash: txHash}, nil
}

// GetTransactionStatus
func (s *grpcServer) GetTransactionStatus(ctx context.Context, in *pb.GetTransactionStatusRequest) (*pb.GetTransactionStatusReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetTransactionStatus %s", in.TxHash)
	status, err := s.nodeClient.GetTxStatusByTxID(ctx, in.TxHash)
	if err != nil {
		log.Errorf("send transaction fails: %s", err)
		return nil, err
	}
	return &pb.GetTransactionStatusReply{Status: string(status)}, nil
}
