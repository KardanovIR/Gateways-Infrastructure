package server

import (
	"context"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/server/converter"
)

// Get account's next nonce
func (s *grpcServer) GetNextNonce(ctx context.Context, in *pb.AddressRequest) (*pb.GetNextNonceReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetNextNonce for address %s", in.Address)

	nonce, err := s.nodeClient.GetNextNonce(ctx, in.Address)
	if err != nil {
		log.Errorf(" getting nonce fails: %s", err)
		return nil, err
	}

	return &pb.GetNextNonceReply{Nonce: int64(nonce)}, nil
}

// Get account's balance
func (s *grpcServer) GetEthBalance(ctx context.Context, in *pb.AddressRequest) (*pb.GetEthBalanceReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetBalance for address %s", in.Address)

	balance, err := s.nodeClient.GetEthBalance(ctx, in.Address)
	if err != nil {
		log.Errorf(" getting balance fails: %s", err)
		return nil, err
	}

	return &pb.GetEthBalanceReply{Balance: converter.ToTargetAmountStr(balance)}, nil
}

// Get account's balance and balances of requested tokens// Get account's balance
func (s *grpcServer) GetAllBalance(ctx context.Context, in *pb.GetAllBalanceRequest) (*pb.GetAllBalanceReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetBalanceIncludedTokens for address %s, balance for contracts %v", in.Address, in.Contracts)

	balance, err := s.nodeClient.GetAllBalances(ctx, in.Address, in.Contracts...)
	if err != nil {
		log.Errorf("getting token's balance fails: %s", err)
		return nil, err
	}
	tokenBalances := make([]*pb.GetAllBalanceReply_TokenBalance, 0, len(balance.Tokens))
	for c, amount := range balance.Tokens {
		tokenBalances = append(tokenBalances,
			&pb.GetAllBalanceReply_TokenBalance{Contract: c, Amount: converter.ToTargetAmountStr(amount)},
		)
	}
	return &pb.GetAllBalanceReply{
		Amount:        converter.ToTargetAmountStr(balance.Amount),
		TokenBalances: tokenBalances,
	}, nil
}
