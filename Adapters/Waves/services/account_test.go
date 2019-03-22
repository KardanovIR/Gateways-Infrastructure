package services

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestNodeClient_GetAllBalance(t *testing.T) {
	ctx, log := beforeTest()
	address := "3N7DGmkCmUgMo9jpuUekUhCMpiBRR1Zm51p"
	wBtcsAssetId := "DWgwcZTMhSvnyYCoWLRUXXSH1RSkzThXLJhww9gwkqdn"
	balance, err := GetNodeClient().GetAllBalances(ctx, address)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, balance.Amount, uint64(190000000))
	assert.Equal(t, len(balance.Assets), 1)
	assert.Equal(t, balance.Assets[wBtcsAssetId], uint64(101000))
}
