package server

import (
	"context"
	"strconv"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
)

// create raw transaction by senders's public key
func (s *grpcServer) GetRawTransaction(ctx context.Context, in *pb.RawTransactionRequest) (
	*pb.RawTransactionReply, error) {

	log := logger.FromContext(ctx)
	log.Infof("GetRawTransaction: address %s to %s send %s ergo", in.AddressFrom, in.AddressTo, in.Amount)
	amount, err := strconv.Atoi(in.Amount)
	if err != nil {
		return nil, err
	}
	tx, err := s.nodeClient.CreateRawTxBySendersAddress(ctx, in.AddressFrom, in.AddressTo, uint64(amount))
	if err != nil {
		log.Errorf("get raw transaction fails: %s", err)
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
