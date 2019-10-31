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

	return &pb.GetNextNonceReply{Nonce: nonce}, nil
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
	b, err := s.converter.ToTargetAmountStr(ctx, balance, "")
	if err != nil {
		log.Errorf("convert balance fails: %s", err)
		return nil, err
	}
	return &pb.GetEthBalanceReply{Balance: b}, nil
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
		am, err := s.converter.ToTargetAmountStr(ctx, amount, c)
		if err != nil {
			log.Errorf("convert token's %s balance %s fails: %s", c, amount, err)
			return nil, err
		}
		tokenBalances = append(tokenBalances,
			&pb.GetAllBalanceReply_TokenBalance{Contract: c, Amount: am},
		)
	}
	ethAmount, err := s.converter.ToTargetAmountStr(ctx, balance.Amount, "")
	if err != nil {
		log.Errorf("convert token's balances fails: %s", err)
		return nil, err
	}
	return &pb.GetAllBalanceReply{
		Amount:        ethAmount,
		TokenBalances: tokenBalances,
	}, nil
}

func (s *grpcServer) AllowanceAmountForAddress(ctx context.Context, in *pb.AllowanceAmountForAddressRequest) (*pb.AllowanceAmountForAddressReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("AllowanceAmountForAddress for owner address %s for address %s (contract %s)", in.OwnerAddress,
		in.SenderAddress, in.Contract)

	amount, err := s.nodeClient.GetErc20AllowanceAmount(ctx, in.OwnerAddress, in.Contract, in.SenderAddress)
	if err != nil {
		log.Errorf("getting allowance amount for %s of owner's address %s failed: %s", in.SenderAddress, in.OwnerAddress, err)
		return nil, err
	}
	convertedAm, err := s.converter.ToTargetAmountStr(ctx, amount, in.Contract)
	if err != nil {
		log.Errorf("convert token's %s balances %s fails: %s", in.Contract, amount, err)
		return nil, err
	}
	return &pb.AllowanceAmountForAddressReply{
		Amount: convertedAm,
	}, nil
}
