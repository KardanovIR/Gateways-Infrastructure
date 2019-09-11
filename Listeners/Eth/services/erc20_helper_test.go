package services

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestNodeClient_ParseERC20TransferParams(t *testing.T) {
	// data from tx 0xfa4c3cf5cf4578a5b051039db0b20061471fafd09bfaa388f93b74b79f03f372 transferFrom in hex format
	// transferFrom tx
	txData := "0x23b872dd000000000000000000000000365950a59f653fd501b8627c2118a93e7fd8e062000000000000000000000000eba28b35e8a02cf648fe1d7d0a767676ee2069b5000000000000000000000000000000000000000000000000058d15e176280000"
	bytes, _ := hexutil.Decode(txData)
	transfer, err := ParseERC20TransferParams(bytes)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, "400000000000000000", transfer.Value.String())
	assert.Equal(t, "0xeBa28b35e8A02Cf648fe1d7d0A767676ee2069B5", transfer.To.String())
	assert.Equal(t, "0x365950a59F653Fd501B8627c2118a93E7fD8e062", transfer.From.String())
	ok, _ := CheckERC20Transfers(bytes)
	assert.True(t, ok)
}
