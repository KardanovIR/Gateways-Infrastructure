package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"go.mongodb.org/mongo-driver/bson"
)

func TestTaskRepository(t *testing.T) {
	ctx, log := beforeTest()
	dbConf := config.DB{Host: "localhost:27017", Name: "unspentTxDb_test"}
	rep, err := New(ctx, dbConf)
	if err != nil {
		log.Fatal("Can't create repository: ", err)
	}
	err = rep.AddUnspentTx(ctx, models.UnspentTx{Address: "2NB7EbQxbut764gJebZoJeAi8grzVrEVfPa",
		TxHash: "c51b19369b38cbdcfa3b902cf565cc00b3ae901479902b5262b1dfcacca295df", Amount: 202, OutputN: 1})
	assert.Nil(t, err)
	err = rep.AddUnspentTx(ctx, models.UnspentTx{Address: "2NB7EbQxbut764gJebZoJeAi8grzVrEVfPa",
		TxHash: "184d5060a1c0cde6213284db7464eb0a866ccbc40d6dd7acecacee38ad84982a", Amount: 100, OutputN: 1})
	assert.Nil(t, err)
	err = rep.AddUnspentTx(ctx, models.UnspentTx{Address: "mxRx7FfxFs1xpihZqQNzVE2c9yoAM1Wd7X",
		TxHash: "c51b19369b38cbdcfa3b902cf565cc00b3ae901479902b5262b1dfcacca295df", Amount: 3, OutputN: 0})
	assert.Nil(t, err)
	count, err := rep.(*repository).unspentTxC.CountDocuments(ctx, bson.D{{}})
	assert.Nil(t, err)
	assert.Equal(t, int64(3), count)
	list, err := rep.GetUnspentTxListForAddress(ctx, "2NB7EbQxbut764gJebZoJeAi8grzVrEVfPa")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	list2, err := rep.GetUnspentTxListForAddress(ctx, "mxRx7FfxFs1xpihZqQNzVE2c9yoAM1Wd7X")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list2))
	assert.Equal(t, uint32(0), list2[0].OutputN)
	assert.Equal(t, uint64(3), list2[0].Amount)
	assert.Equal(t, "c51b19369b38cbdcfa3b902cf565cc00b3ae901479902b5262b1dfcacca295df", list2[0].TxHash)
	u, err := rep.GetUnspentTxByTxHashAndOutputNumber(ctx, "c51b19369b38cbdcfa3b902cf565cc00b3ae901479902b5262b1dfcacca295df", 1)
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), u.OutputN)
	assert.Equal(t, uint64(202), u.Amount)
	assert.Equal(t, "2NB7EbQxbut764gJebZoJeAi8grzVrEVfPa", u.Address)
	assert.Equal(t, "c51b19369b38cbdcfa3b902cf565cc00b3ae901479902b5262b1dfcacca295df", u.TxHash)
	// unknown input
	u2, err := rep.GetUnspentTxByTxHashAndOutputNumber(ctx, "c51b19369b38cbdcfa3b902cf565cc00b3ae901479902b5262b1dfcacca295df", 3)
	assert.Nil(t, err)
	assert.Nil(t, u2)

	du, err := rep.DeleteUnspentTx(ctx, "c51b19369b38cbdcfa3b902cf565cc00b3ae901479902b5262b1dfcacca295df", 3)
	assert.Nil(t, err)
	assert.Nil(t, du)

	du2, err := rep.DeleteUnspentTx(ctx, "c51b19369b38cbdcfa3b902cf565cc00b3ae901479902b5262b1dfcacca295df", 1)
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), du2.OutputN)
	assert.Equal(t, uint64(202), du2.Amount)
	assert.Equal(t, "2NB7EbQxbut764gJebZoJeAi8grzVrEVfPa", du2.Address)
	assert.Equal(t, "c51b19369b38cbdcfa3b902cf565cc00b3ae901479902b5262b1dfcacca295df", du2.TxHash)
	du3, err := rep.DeleteUnspentTx(ctx, "c51b19369b38cbdcfa3b902cf565cc00b3ae901479902b5262b1dfcacca295df", 0)
	assert.Nil(t, err)
	assert.NotNil(t, du3)
	du4, err := rep.DeleteUnspentTx(ctx, "184d5060a1c0cde6213284db7464eb0a866ccbc40d6dd7acecacee38ad84982a", 1)
	assert.Nil(t, err)
	assert.NotNil(t, du4)
	count2, err := rep.(*repository).unspentTxC.CountDocuments(ctx, bson.D{{}})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), count2)
}

func beforeTest() (ctx context.Context, log logger.ILogger) {
	ctx = context.Background()
	log, _ = logger.Init(false, logger.DEBUG)
	return
}
