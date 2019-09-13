package server

import (
	"context"
	"fmt"
	"strconv"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
)

// create raw transaction with one recipient
func (s *grpcServer) GetRawTransaction(ctx context.Context, in *pb.RawTransactionRequest) (
	*pb.RawTransactionReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetRawTransaction: address %s to %s send %s ergo", in.AddressFrom, in.AddressTo, in.Amount)
	outputs := make([]*models.Output, 1)
	amount, err := strconv.ParseUint(in.Amount, 10, 64)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	outputs[0] = &models.Output{Address: in.AddressTo, Amount: amount}
	tx, err := s.nodeClient.CreateRawTx(ctx, in.AddressFrom, outputs)
	if err != nil {
		log.Errorf("get raw transaction fails: %s", err)
		return nil, err
	}
	return &pb.RawTransactionReply{Tx: tx}, nil
}

// create raw transaction with many recipients
func (s *grpcServer) GetRawMassTransaction(ctx context.Context, in *pb.RawMassTransactionRequest) (*pb.RawTransactionReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetRawMassTransaction: address %s to %v ergo", in.AddressFrom, in.Outs)
	if len(in.AddressFrom) == 0 || len(in.Outs) == 0 {
		return nil, fmt.Errorf("wrong parameters %s, %+v", in.AddressFrom, in.Outs)
	}
	outputs := make([]*models.Output, len(in.Outs))
	for i, o := range in.Outs {
		amount, err := strconv.ParseUint(o.Amount, 10, 64)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		outputs[i] = &models.Output{Address: o.AddressTo, Amount: amount}
	}
	tx, err := s.nodeClient.CreateRawTx(ctx, in.AddressFrom, outputs)
	if err != nil {
		log.Errorf("GetRawMassTransaction fails: %s", err)
		return nil, err
	}
	return &pb.RawTransactionReply{Tx: tx}, nil
}

// send transaction
func (s *grpcServer) SendTransaction(ctx context.Context, in *pb.SendTransactionRequest) (
	*pb.SendTransactionReply, error) {

	log := logger.FromContext(ctx)
	log.Info("SendTransaction")

	txId, err := s.nodeClient.SendTransaction(ctx, in.Tx)
	if err != nil {
		log.Errorf("sending of transaction fails: %s", err)
		return nil, err
	}
	return &pb.SendTransactionReply{TxId: txId}, nil
}

// get transaction by hash
func (s *grpcServer) TransactionByHash(ctx context.Context, in *pb.TransactionByHashRequest) (*pb.TransactionByHashReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("TransactionByHash for %s", in.TxHash)
	tx, err := s.nodeClient.TransactionByHash(ctx, in.TxHash)
	if err != nil {
		log.Errorf("getting transaction by hash fails: %s", err)
		return nil, err
	}
	inputs := make([]*pb.InputOutput, 0)
	for _, in := range tx.Inputs {
		inputs = append(inputs, &pb.InputOutput{Amount: in.Amount, Address: in.Address})
	}
	outputs := make([]*pb.InputOutput, 0)
	for _, out := range tx.Outputs {
		outputs = append(outputs, &pb.InputOutput{Amount: out.Amount, Address: out.Address})
	}
	return &pb.TransactionByHashReply{
		SenderAddress:    tx.From,
		SendersPublicKey: tx.SenderPublicKey,
		RecipientAddress: tx.To,
		Amount:           tx.Amount,
		Fee:              tx.Fee,
		AssetId:          tx.AssetId,
		Status:           string(tx.Status),
		TxHash:           tx.TxHash,
		Data:             tx.Data,
		Inputs:           inputs,
		Outputs:          outputs,
	}, nil
}
