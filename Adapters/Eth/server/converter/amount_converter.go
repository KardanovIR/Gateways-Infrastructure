package converter

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/services"
)

const (
	nodeDecimals = 18
)

type IConverter interface {
	ToNodeAmount(ctx context.Context, a *big.Int, contract string) (*big.Int, error)
	ToTargetAmountStr(ctx context.Context, a *big.Int, contract string) (string, error)
	ToTargetAmount(ctx context.Context, a *big.Int, contract string) (*big.Int, error)
	ToCommissionStr(a *big.Int) string
}

type converter struct {
	sync.RWMutex
	maxTargetDecimals   int64
	multiplierToNode    *big.Int
	contractsMultiplier map[string]*big.Int
	contractProvider    services.IDecimalsContractProvider
}

// method should be call one time
func Init(ctx context.Context, maxTargetDecimals int64, contractProvider services.IDecimalsContractProvider) IConverter {
	if nodeDecimals < maxTargetDecimals {
		logger.FromContext(ctx).Fatalf("wrong parameter 'maxTargetDecimals' = %d. It should be less than %d", maxTargetDecimals, nodeDecimals)
	}
	c := converter{
		maxTargetDecimals:   maxTargetDecimals,
		multiplierToNode:    countMultiplier(nodeDecimals, maxTargetDecimals),
		contractsMultiplier: make(map[string]*big.Int),
		contractProvider:    contractProvider,
	}
	return &c
}

func (c *converter) getMultiplierForContract(ctx context.Context, contract string) (*big.Int, error) {
	if len(contract) == 0 {
		return c.multiplierToNode, nil
	}
	if m, ok := c.readFromContractsMultiplier(contract); ok {
		return m, nil
	}
	decimals, err := c.contractProvider.Decimals(ctx, contract)
	if err != nil {
		return nil, err
	}
	if decimals == 0 {
		logger.FromContext(ctx).Warnf("request for decimals for contract %s return 0!", contract)
	}
	if decimals > nodeDecimals {
		err := fmt.Errorf("decimals %d > node's decimal %d for contract %s", decimals, nodeDecimals, contract)
		logger.FromContext(ctx).Error(err)
		return nil, err
	}
	multiplier := countMultiplier(decimals, c.maxTargetDecimals)
	c.Lock()
	defer c.Unlock()
	c.contractsMultiplier[contract] = multiplier
	return multiplier, nil
}

func (c *converter) readFromContractsMultiplier(contract string) (*big.Int, bool) {
	c.RLock()
	defer c.RUnlock()
	m, ok := c.contractsMultiplier[contract]
	return m, ok
}

func countMultiplier(current, maxTarget int64) *big.Int {
	// if current is not more than max -> not need to convert it
	if current <= maxTarget {
		big.NewInt(1)
	}
	diff := current - maxTarget
	multiplier := int64(1)
	for i := int64(0); i < diff; i++ {
		multiplier *= 10
	}
	return big.NewInt(multiplier)
}

// ToNodeAmount return amount with decimals used in eth blockchain
func (c *converter) ToNodeAmount(ctx context.Context, a *big.Int, contract string) (*big.Int, error) {
	if a == nil {
		return new(big.Int), nil
	}
	multiplier, err := c.getMultiplierForContract(ctx, contract)
	if err != nil {
		return nil, err
	}
	return new(big.Int).Mul(a, multiplier), nil
}

func (c *converter) ToTargetAmountStr(ctx context.Context, a *big.Int, contract string) (string, error) {
	if a == nil {
		return "0", nil
	}
	convertedAmount, err := c.ToTargetAmount(ctx, a, contract)
	if err != nil {
		return "", err
	}
	return convertedAmount.String(), nil
}

func (c *converter) ToTargetAmount(ctx context.Context, a *big.Int, contract string) (*big.Int, error) {
	if a == nil {
		return new(big.Int), nil
	}
	multiplier, err := c.getMultiplierForContract(ctx, contract)
	if err != nil {
		return nil, err
	}
	return new(big.Int).Div(a, multiplier), nil
}

func (c *converter) ToCommissionStr(a *big.Int) string {
	if a == nil {
		return "0"
	}
	// fee is paid in eth -> not need to ask for decimals - we know it
	targetAmount := new(big.Int).Div(a, c.multiplierToNode)
	rounded := new(big.Int).Add(targetAmount, big.NewInt(1))
	return rounded.String()
}
