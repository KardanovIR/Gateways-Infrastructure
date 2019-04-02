package repositories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
)

func TestRepository_PutChainState(t *testing.T) {
	ctx, log := beforeTest()
	r := GetRepository().(*repository)
	defer func() {
		// drop collection after test complete successful
		if !t.Failed() {
			if err := r.chainStateC.Drop(ctx); err != nil {
				log.Error(err)
			}
		}
	}()
	// get chain state first time (empty state)
	cs, err := GetRepository().GetLastChainState(ctx, models.Ethereum)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Nil(t, cs)
	// put state (insert)
	st1 := models.ChainState{Timestamp: time.Now(), ChainType: models.Ethereum, LastBlock: 123456}
	st2, err := GetRepository().PutChainState(ctx, &st1)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, int64(123456), st2.LastBlock)

	//  put state (update)
	st3 := models.ChainState{Timestamp: time.Now(), ChainType: models.Ethereum, LastBlock: 200000}
	st4, err := GetRepository().PutChainState(ctx, &st3)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, int64(200000), st4.LastBlock)
	// get chain state after insert and update
	cs2, err := GetRepository().GetLastChainState(ctx, models.Ethereum)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.NotNil(t, cs2)
	if cs2 == nil {
		t.FailNow()
	}
	assert.Equal(t, int64(200000), cs2.LastBlock)
}
