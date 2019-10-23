package server

import (
	"context"
	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
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
	log.Infof("GetUnspentInputs for address %s", in.Address)
	var utxo, err = s.dataClient.GetUnspentInputs(ctx, in.Address)
	if err != nil {
		log.Errorf("getting all unspent inputs fails: %s", err)
		return nil, err
	}
	outputs := make([]*pb.UnspentInput, 0)
	for _, input := range utxo {
		outputs = append(outputs, &pb.UnspentInput{
			TxId:          input.TxId,
			Address:       input.Address,
			Vout:          uint32(input.Vout),
			ScriptPubKey:  input.ScriptPubKey,
			Amount:        input.Amount,
			Satoshis:      uint32(input.Satoshis),
			Height:        uint32(input.Height),
			Confirmations: uint32(input.Confirmations)})
	}
	return &pb.GetUnspentInputsReply{Utxo: outputs}, nil
}
