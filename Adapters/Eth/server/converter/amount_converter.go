package converter

import (
	"context"
	"math/big"

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
	m := countMultiplier(nodeDecimals, maxTargetDecimals)
	c := converter{
		multiplierToNode:    m,
		contractsMultiplier: make(map[string]*big.Int),
		contractProvider:    contractProvider,
		maxTargetDecimals:   maxTargetDecimals,
	}
	return &c
}

func (c *converter) getMultiplierForContract(ctx context.Context, contract string) (*big.Int, error) {
	if len(contract) == 0 {
		return c.multiplierToNode, nil
	}
	if m, ok := c.contractsMultiplier[contract]; ok {
		return m, nil
	}
	decimals, err := c.contractProvider.Decimals(ctx, contract)
	if err != nil {
		return nil, err
	}
	if decimals == 0 {
		logger.FromContext(ctx).Warnf("request for decimals for contract %s return 0!", contract)
	}
	m := countMultiplier(decimals, c.maxTargetDecimals)
	c.contractsMultiplier[contract] = m
	return m, nil
}

func countMultiplier(current, maxTarget int64) *big.Int {
	// if current is not more than max -> not need to convert it
	if current <= maxTarget {
		big.NewInt(1)
	}
	diff := current - maxTarget
	m := int64(1)
	for i := int64(0); i < diff; i++ {
		m *= 10
	}
	return big.NewInt(m)
}

// ToNodeAmount return amount with decimals used in eth blockchain
func (c *converter) ToNodeAmount(ctx context.Context, a *big.Int, contract string) (*big.Int, error) {
	if a == nil {
		return new(big.Int), nil
	}
	m, err := c.getMultiplierForContract(ctx, contract)
	if err != nil {
		return nil, err
	}
	return new(big.Int).Mul(a, m), nil
}

func (c *converter) ToTargetAmountStr(ctx context.Context, a *big.Int, contract string) (string, error) {
	if a == nil {
		return "0", nil
	}
	am, err := c.ToTargetAmount(ctx, a, contract)
	if err != nil {
		return "", err
	}
	return am.String(), nil
}

func (c *converter) ToTargetAmount(ctx context.Context, a *big.Int, contract string) (*big.Int, error) {
	if a == nil {
		return new(big.Int), nil
	}
	m, err := c.getMultiplierForContract(ctx, contract)
	if err != nil {
		return nil, err
	}
	return new(big.Int).Div(a, m), nil
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
