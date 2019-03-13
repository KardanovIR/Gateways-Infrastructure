package services

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// GetBalance return current balance for address
func (cl *nodeClient) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	return cl.ethClient.BalanceAt(ctx, common.HexToAddress(address), nil)
}

func (cl *nodeClient) GetNextNonce(ctx context.Context, address string) (uint64, error) {
	a := common.HexToAddress(address)
	return cl.ethClient.PendingNonceAt(ctx, a)
}
