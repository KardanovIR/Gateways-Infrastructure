package services

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/models"
)

// GetEthBalance return current balance for address
func (cl *nodeClient) GetEthBalance(ctx context.Context, address string) (*big.Int, error) {
	return cl.ethClient.BalanceAt(ctx, common.HexToAddress(address), nil)
}

func (cl *nodeClient) GetNextNonce(ctx context.Context, address string) (uint64, error) {
	a := common.HexToAddress(address)
	return cl.ethClient.PendingNonceAt(ctx, a)
}

// GetAllBalances return eth balance and token's balances for address.
// Token's balances requested only for tokens which contract was passed to parameters
func (cl *nodeClient) GetAllBalances(ctx context.Context, address string, contracts ...string) (*models.AccountBalance, error) {
	log := logger.FromContext(ctx)
	log.Debugf("Call method 'GetAllBalances' for %s", address)
	balances := models.AccountBalance{Tokens: make(map[string]*big.Int)}
	ethBalance, err := cl.GetEthBalance(ctx, address)
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

func (cl *nodeClient) GetErc20AllowanceAmount(ctx context.Context, ownerAddress string, contractAddress string,
	senderAddress string) (*big.Int, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'CreateErc20Allowance': check sending allowance for sender %s for address %s; contract %s",
		senderAddress, ownerAddress, contractAddress)

	sender := common.HexToAddress(senderAddress)
	owner := common.HexToAddress(ownerAddress)
	data, err := ERC20AllowanceForSender(owner, sender)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	contract := common.HexToAddress(contractAddress)
	msg := ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}
	contractResponse, err := cl.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		log.Error("can't get allowance for %s (contract %s) for address %s: %s", sender, contract, ownerAddress, err)
		return nil, err
	}
	return new(big.Int).SetBytes(contractResponse), nil
}
