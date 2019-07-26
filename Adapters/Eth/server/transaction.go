package server

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/server/converter"
)

// create raw transaction
func (s *grpcServer) GetRawTransaction(ctx context.Context, in *pb.RawTransactionRequest) (*pb.RawTransactionReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("getRawTransaction %+v", in)
	amount, ok := new(big.Int).SetString(in.Amount, 10)
	if !ok {
		err := fmt.Errorf("wrong amount value: %s", in.Amount)
		log.Error(err)
		return nil, err
	}
	amount = converter.ToNodeAmount(amount)
	var tx []byte
	var err error
	if len(in.Contract) > 0 {
		tx, err = s.nodeClient.CreateErc20TokensRawTransaction(ctx, in.AddressFrom, in.Contract, in.AddressTo, amount)
	} else {
		tx, err = s.nodeClient.CreateRawTransaction(ctx, in.AddressFrom, in.AddressTo, amount)
	}
	if err != nil {
		log.Errorf("transaction's creation fails: %s", err)
		return nil, err
	}
	return &pb.RawTransactionReply{Tx: tx}, nil
}

// create raw transaction to transfer erc-20 tokens
func (s *grpcServer) GetErc20RawTransaction(ctx context.Context, in *pb.Erc20RawTransactionRequest) (*pb.RawTransactionReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetErc20RawTransaction %+v", in)
	amount, ok := new(big.Int).SetString(in.Amount, 10)
	if !ok {
		err := fmt.Errorf("wrong amount value: %s", in.Amount)
		log.Error(err)
		return nil, err
	}
	amount = converter.ToNodeAmount(amount)
	var tx, err = s.nodeClient.CreateErc20TokensRawTransaction(ctx, in.AddressFrom, in.Contract, in.AddressTo, amount)
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
		log.Errorf("get transaction status fails: %s", err)
		return nil, err
	}
	return &pb.GetTransactionStatusReply{Status: string(status)}, nil
}

// TransactionByHash
func (s *grpcServer) TransactionByHash(ctx context.Context, in *pb.TransactionByHashRequest) (*pb.TransactionByHashReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("TransactionByHash %s", in.TxHash)
	tx, err := s.nodeClient.TransactionInfo(ctx, in.TxHash)
	if err != nil {
		log.Errorf("get transaction by hash fails: %s", err)
		return nil, err
	}
	return &pb.TransactionByHashReply{
		SenderAddress:    tx.From,
		RecipientAddress: tx.To,
		Amount:           converter.ToTargetAmountStr(tx.Amount),
		Fee:              converter.ToTargetAmountStr(tx.Fee),
		AssetId:          tx.Contract,
		Status:           string(tx.Status),
		TxHash:           tx.TxHash,
		Data:             tx.Data,
	}, nil
}
