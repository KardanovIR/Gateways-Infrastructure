package converter

import (
	"context"
	"strconv"

	"github.com/shopspring/decimal"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
)

func GetIntFromFloat(ctx context.Context, value float64) (uint64, error) {
	amount := decimal.NewFromFloat(value).Shift(8)
	amount.Truncate(0)
	strInt := amount.String()
	result, err := strconv.ParseUint(strInt, 10, 64)
	if err != nil {
		logger.FromContext(ctx).Errorf("convert amount from float64 %f to uint64 fails %s", value, err)
		return 0, err
	}
	return result, nil
}
