package converter

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToTargetAmount(t *testing.T) {
	Init(context.Background(), 8)
	a, _ := new(big.Int).SetString("123456", 10)
	roundedA, _ := new(big.Int).SetString("123457", 10)
	aNode, _ := new(big.Int).SetString("1234560000000000", 10)
	assert.Equal(t, aNode.String(), ToNodeAmount(a).String())
	assert.Equal(t, a.String(), ToTargetAmountStr(aNode))
	assert.Equal(t, roundedA.String(), ToCommissionStr(aNode))

	aNode2, _ := new(big.Int).SetString("123456789087654321", 10)
	assert.Equal(t, "12345678", ToTargetAmountStr(aNode2))
}
