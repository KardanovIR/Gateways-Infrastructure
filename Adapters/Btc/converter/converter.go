package converter

import (
	"github.com/shopspring/decimal"
	"strconv"
)

func GetIntFromFloat(value float64) (uint64, error) {
	amount := decimal.NewFromFloat(value).Shift(8)
	amount.Truncate(0)
	strInt := amount.String()
	result, err := strconv.ParseUint(strInt, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}
