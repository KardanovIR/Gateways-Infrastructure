package services

import (
	"context"
	"math/big"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
)

func (cl *nodeClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'SuggestGasPrice'")
	return cl.ethClient.SuggestGasPrice(ctx)
}

func (cl *nodeClient) SuggestFee(ctx context.Context) (*big.Int, error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'SuggestFee'")
	gasPrice, err := cl.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	gasLimit := new(big.Int).SetInt64(gasLimitForMoneyTransfer)
	return new(big.Int).Mul(gasPrice, gasLimit), nil
}
