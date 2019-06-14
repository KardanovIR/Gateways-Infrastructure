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
	log.Infof("GetRawTransactionBySendersPublicKey: pk %s to %s send %s (asset %s)", in.SendersPublicKey,
		in.AddressTo, in.Amount, in.AssetId)
	amount, err := strconv.Atoi(in.Amount)
	if err != nil {
		return nil, err
	}
	tx, err := s.nodeClient.CreateRawTxBySendersPublicKey(ctx, in.SendersPublicKey, in.AddressTo, uint64(amount), in.AssetId)
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
	}, nil
}
