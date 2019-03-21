package services

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/models"
)

// GetBalance return current balance for address
func (cl *nodeClient) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	return cl.ethClient.BalanceAt(ctx, common.HexToAddress(address), nil)
}

func (cl *nodeClient) GetNextNonce(ctx context.Context, address string) (uint64, error) {
	a := common.HexToAddress(address)
	return cl.ethClient.PendingNonceAt(ctx, a)
}

// GetBalance return eth balance for address and token's balance for contracts
func (cl *nodeClient) GetTokenBalance(ctx context.Context, address string, contracts ...string) (*models.AccountBalance, error) {
	log := logger.FromContext(ctx)
	log.Debugf("Call method 'GetTokenBalance' for %s", address)
	balances := models.AccountBalance{Tokens: make(map[string]*big.Int)}
	ethBalance, err := cl.GetBalance(ctx, address)
	if err != nil {
		return nil, err
	}
	balances.Amount = ethBalance
	if len(contracts) == 0 {
		return &balances, nil
	}
	for _, c := range contracts {
		tokenBalance, err := cl.contractProvider.BalanceOf(ctx, address, c)
		if err != nil {
			return nil, err
		}
		balances.Tokens[c] = tokenBalance
		log.Debugf("token balance for contract %s for account %s: %s", c, address, tokenBalance)
	}
	return &balances, nil
}
