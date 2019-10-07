package converter

import (
	"context"
	"strconv"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
)

const (
	btcDecimals             = 8
	decimalBaseToFormatUint = 10
)

var multiplierToNode uint64

// method should be call one time
func Init(ctx context.Context, targetDecimals int) {
	if btcDecimals < targetDecimals {
		logger.FromContext(ctx).Fatalf("wrong parameter 'targetDecimals' = %d. It should be less than %d", targetDecimals, btcDecimals)
	}
	diff := btcDecimals - targetDecimals
	multiplierToNode = uint64(1)
	for i := 0; i < diff; i++ {
		multiplierToNode *= 10
	}
}

func ToNodeAmount(a uint64) uint64 {
	return a * multiplierToNode
}

func ToTargetAmount(a float64) uint64 {
	return uint64(a *btcDecimals)
}

func ToTargetAmountStr(a float64) string {
	return strconv.FormatUint(ToTargetAmount(a), decimalBaseToFormatUint)
}
