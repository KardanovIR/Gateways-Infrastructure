package converter

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/services"
)

func TestToTargetAmount(t *testing.T) {
	c := Init(context.Background(), 8, nil)
	ctx := context.Background()

	a, _ := new(big.Int).SetString("123456", 10)
	roundedA, _ := new(big.Int).SetString("123457", 10)
	aNode, _ := new(big.Int).SetString("1234560000000000", 10)
	r, err := c.ToNodeAmount(ctx, a, "")
	assert.Nil(t, err)
	assert.Equal(t, aNode.String(), r.String())
	r2, err := c.ToTargetAmountStr(ctx, aNode, "")
	assert.Nil(t, err)
	assert.Equal(t, a.String(), r2)
	assert.Equal(t, roundedA.String(), c.ToCommissionStr(aNode))

	aNode2, _ := new(big.Int).SetString("123456789087654321", 10)
	r3, err := c.ToTargetAmountStr(ctx, aNode2, "")
	assert.Nil(t, err)
	assert.Equal(t, "12345678", r3)
}

func beforeTest() (context.Context, logger.ILogger) {
	ctx := context.Background()
	log, _ := logger.Init(false, logger.DEBUG)
	err := config.Load("./testdata/config_test.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = services.New(ctx, config.Cfg.Node)
	if err != nil {
		log.Fatal(err)
	}
	return ctx, log
}

func TestNodeClient_CallDecimalsInContract(t *testing.T) {
	ctx, _ := beforeTest()
	c := Init(ctx, 8, services.GetNodeClient().GetContractProvider())
	// decimals = 18 (eth)
	r, err := c.ToTargetAmountStr(ctx, big.NewInt(1220000000000), "")
	assert.Nil(t, err)
	assert.Equal(t, "122", r)
	na, err := c.ToNodeAmount(ctx, big.NewInt(122), "")
	assert.Nil(t, err)
	assert.Equal(t, "1220000000000", na.String())

	// decimals = 18
	r2, err := c.ToTargetAmountStr(ctx, big.NewInt(1220000000000), "0x722dd3F80BAC40c951b51BdD28Dd19d435762180")
	assert.Nil(t, err)
	assert.Equal(t, "122", r2)
	na2, err := c.ToNodeAmount(ctx, big.NewInt(122), "0x722dd3F80BAC40c951b51BdD28Dd19d435762180")
	assert.Nil(t, err)
	assert.Equal(t, "1220000000000", na2.String())

	// decimals = 6
	r3, err := c.ToTargetAmountStr(ctx, big.NewInt(1220000000000), "0xa33e09ec9f0aaecaad4887eafe6f9f1e5d1e812d")
	assert.Nil(t, err)
	assert.Equal(t, "1220000000000", r3)
	na3, err := c.ToNodeAmount(ctx, big.NewInt(122), "0xa33e09ec9f0aaecaad4887eafe6f9f1e5d1e812d")
	assert.Nil(t, err)
	assert.Equal(t, "122", na3.String())
}
