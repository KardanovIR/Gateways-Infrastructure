package repositories

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
)

func TestRepository_PutChainState(t *testing.T) {
	ctx, log := beforeTest()
	rep, err := connect(ctx, "localhost:27017", "taskDb_test")
	if err != nil {
		log.Fatal("Can't create db connection: ", err)
	}
	r := rep.(*repository)
	defer func() {
		// drop chain state after test complete successful
		if _, err := r.chainStateC.DeleteOne(ctx, bson.D{{"chaintype", models.Ethereum}}); err != nil {
			log.Error(err)
		}
	}()
	// get chain state first time (empty state)
	cs, err := rep.GetLastChainState(ctx, models.Ethereum)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Nil(t, cs)
	// put state (insert)
	st1 := models.ChainState{Timestamp: time.Now(), ChainType: models.Ethereum, LastBlock: 123456}
	st2, err := rep.PutChainState(ctx, &st1)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, int64(123456), st2.LastBlock)

	//  put state (update)
	st3 := models.ChainState{Timestamp: time.Now(), ChainType: models.Ethereum, LastBlock: 200000}
	st4, err := rep.PutChainState(ctx, &st3)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, int64(200000), st4.LastBlock)
	// get chain state after insert and update
	cs2, err := rep.GetLastChainState(ctx, models.Ethereum)
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
