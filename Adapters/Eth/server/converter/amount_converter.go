package converter

import (
	"context"
	"math/big"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
)

const (
	nodeDecimals = 18
)

var multiplierToNode *big.Int

// method should be call one time
func Init(ctx context.Context, targetDecimals int) {
	if nodeDecimals < targetDecimals {
		logger.FromContext(ctx).Fatalf("wrong parameter 'targetDecimals' = %d. It should be less than %d", targetDecimals, nodeDecimals)
	}
	diff := nodeDecimals - targetDecimals
	m := int64(1)
	for i := 0; i < diff; i++ {
		m *= 10
	}
	multiplierToNode = big.NewInt(m)
}

func ToNodeAmount(a *big.Int) *big.Int {
	return new(big.Int).Mul(a, multiplierToNode)
}

func ToTargetAmountStr(a *big.Int) string {
	if a == nil {
		return "0"
	}
	return ToTargetAmount(a).String()
}

func ToTargetAmount(a *big.Int) *big.Int {
	if a == nil {
		return new(big.Int)
	}
	return new(big.Int).Div(a, multiplierToNode)
}
