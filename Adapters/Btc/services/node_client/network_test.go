package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeClient_FeeRate(t *testing.T) {
	ctx, _ := beforeTest()
	fee, err := GetNodeClient().FeeRateForKByte(ctx)
	assert.Nil(t, err)
	assert.True(t, fee > 0)
}

func TestNodeClient_Fee(t *testing.T) {
	ctx, _ := beforeTest()
	assert.Equal(t, uint64(3161), GetNodeClient().Fee(ctx, 12345, 256))
	assert.Equal(t, uint64(2520), GetNodeClient().Fee(ctx, 10000, 252))
}
