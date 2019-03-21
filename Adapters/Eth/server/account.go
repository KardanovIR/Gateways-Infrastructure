package server

import (
	"context"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
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
func (s *grpcServer) GetBalance(ctx context.Context, in *pb.AddressRequest) (*pb.GetBalanceReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetBalance for address %s", in.Address)

	balance, err := s.nodeClient.GetBalance(ctx, in.Address)
	if err != nil {
		log.Errorf(" getting balance fails: %s", err)
		return nil, err
	}

	return &pb.GetBalanceReply{Balance: balance.String()}, nil
}

// Get account's balance and balances of requested tokens// Get account's balance
func (s *grpcServer) GetBalanceIncludedTokens(ctx context.Context,
	in *pb.GetBalanceIncludedTokensRequest) (*pb.GetBalanceIncludedTokensReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("GetBalanceIncludedTokens for address %s, balance for contracts %v", in.Address, in.Contracts)

	balance, err := s.nodeClient.GetTokenBalance(ctx, in.Address, in.Contracts...)
	if err != nil {
		log.Errorf("getting token's balance fails: %s", err)
		return nil, err
	}
	tokenBalances := make([]*pb.GetBalanceIncludedTokensReply_TokenBalance, 0, len(balance.Tokens))
	for c, amount := range balance.Tokens {
		tokenBalances = append(tokenBalances,
			&pb.GetBalanceIncludedTokensReply_TokenBalance{Contract: c, Amount: amount.String()},
		)
	}
	return &pb.GetBalanceIncludedTokensReply{
		Amount:        balance.Amount.String(),
		TokenBalances: tokenBalances,
	}, nil
}
