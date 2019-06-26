package converter

import (
	"context"
	"strconv"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
)

const (
	ergoDecimals            = 9
	decimalBaseToFormatUint = 10
)

var multiplierToNode uint64

// method should be call one time
func Init(ctx context.Context, targetDecimals int) {
	if ergoDecimals < targetDecimals {
		logger.FromContext(ctx).Fatalf("wrong parameter 'targetDecimals' = %d. It should be less than %d", targetDecimals, ergoDecimals)
	}
	diff := ergoDecimals - targetDecimals
	multiplierToNode = uint64(1)
	for i := 0; i < diff; i++ {
		multiplierToNode *= 10
	}
}

func ToNodeAmount(a uint64) uint64 {
	return a * multiplierToNode
}

func ToTargetAmount(a uint64) uint64 {
	return a / multiplierToNode
}

func ToTargetAmountStr(a uint64) string {
	return strconv.FormatUint(ToTargetAmount(a), decimalBaseToFormatUint)
}
