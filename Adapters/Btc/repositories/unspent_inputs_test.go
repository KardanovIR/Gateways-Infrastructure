package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
)

func TestTaskRepository(t *testing.T) {
	ctx, log := beforeTest()
	dbConf := config.DB{Host: "localhost:27017", Name: "unspent_inputs_test_amount"}
	err := New(ctx, dbConf)
	if err != nil {
		log.Fatal("Can't create repository: ", err)
	}
	defer func() {
		// drop collection after test complete successful
		if err := rep.(*repository).unspentTxC.Drop(ctx); err != nil {
			log.Error(err)
		}
	}()
	unT1 := models.UnspentTx{Address: "2NB7EbQxbut764gJebZoJeAi8grzVrEVfPa",
		TxHash: "c51b19369b38cbdcfa3b902cf565cc00b3ae901479902b5262b1dfcacca295df", Amount: 202, OutputN: 1}
	_, err = rep.(*repository).unspentTxC.InsertOne(ctx, &unT1)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	unT2 := models.UnspentTx{Address: "2NB7EbQxbut764gJebZoJeAi8grzVrEVfPa",
		TxHash: "184d5060a1c0cde6213284db7464eb0a866ccbc40d6dd7acecacee38ad84982a", Amount: 100, OutputN: 1}
	_, err = rep.(*repository).unspentTxC.InsertOne(ctx, &unT2)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	unT3 := models.UnspentTx{Address: "mxRx7FfxFs1xpihZqQNzVE2c9yoAM1Wd7X",
		TxHash: "c51b19369b38cbdcfa3b902cf565cc00b3ae901479902b5262b1dfcacca295df", Amount: 3, OutputN: 1}
	_, err = rep.(*repository).unspentTxC.InsertOne(ctx, &unT3)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	balances, err := GetRepository().GetBalanceForAddresses(ctx, []string{"2NB7EbQxbut764gJebZoJeAi8grzVrEVfPa",
		"mxRx7FfxFs1xpihZqQNzVE2c9yoAM1Wd7X", "mfduw3E2qtTHNyRfkCKH9v9YMdiiAPCwZr"})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(balances))
	assert.Equal(t, "mxRx7FfxFs1xpihZqQNzVE2c9yoAM1Wd7X", balances[0].Address)
	assert.Equal(t, uint64(3), balances[0].Amount)
	assert.Equal(t, "2NB7EbQxbut764gJebZoJeAi8grzVrEVfPa", balances[1].Address)
	assert.Equal(t, uint64(302), balances[1].Amount)
}

func beforeTest() (ctx context.Context, log logger.ILogger) {
	ctx = context.Background()
	log, _ = logger.Init(false, logger.DEBUG)
	return
}
