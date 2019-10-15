package converter

import (
	"math"
	"strconv"
)

const (
	btcDecimals             = 8
	decimalBaseToFormatUint = 10
)

var multiplierToNode uint64


func ToTargetAmount(a float64) uint64 {
	return uint64(a * math.Pow10(btcDecimals))
}

func ToTargetAmountStr(a float64) string {
	return strconv.FormatUint(ToTargetAmount(a), decimalBaseToFormatUint)
}
