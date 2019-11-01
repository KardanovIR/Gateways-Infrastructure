package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeClient_CallDecimalsInContract(t *testing.T) {
	ctx, _ := beforeTest()
	cp := GetNodeClient().(*nodeClient).contractProvider
	d, err := cp.Decimals(ctx, "0x722dd3F80BAC40c951b51BdD28Dd19d435762180")
	assert.Nil(t, err)
	assert.Equal(t, int64(18), d)

	d2, err := cp.Decimals(ctx, "0xa33e09ec9f0aaecaad4887eafe6f9f1e5d1e812d")
	assert.Nil(t, err)
	assert.Equal(t, int64(6), d2)
}
