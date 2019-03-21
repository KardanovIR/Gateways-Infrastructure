package services

import (
	"math/big"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestNodeClient_GetTokenBalance(t *testing.T) {
	ctx, log := beforeTest()
	address := "0xb8ebd916689e773c6657de537151476A3a8259fc"
	// ERC-20 (LINK)
	contract1 := "0x20fe562d797a42dcb3399062ae9546cd06f63280"
	// Test Standard Token (TST)
	contract2 := "0x722dd3F80BAC40c951b51BdD28Dd19d435762180"
	balances, err := GetNodeClient().GetAllBalances(ctx, address, contract1, contract2)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	ethAmount := new(big.Int).SetInt64(4200000000000000)
	assert.Equal(t, balances.Amount.String(), ethAmount.String())
	assert.Equal(t, len(balances.Tokens), 2)
	linkAmount := new(big.Int).SetInt64(325000000000000000)
	assert.Equal(t, balances.Tokens[contract1].String(), linkAmount.String())
	tstAmount, _ := new(big.Int).SetString("19000000000000000000", 10)
	assert.Equal(t, balances.Tokens[contract2].String(), tstAmount.String())
}
