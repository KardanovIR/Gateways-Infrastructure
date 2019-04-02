package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
)

// TestTaskRepository needs to running mongoDB
// check methods PutTask, RemoveTask, FindByAddressOrTxId
func TestTaskRepository(t *testing.T) {
	ctx, log := beforeTest()
	rep, err := connect(ctx, "localhost:27017", "taskDb_test")
	if err != nil {
		log.Fatal("Can't create db connection: ", err)
	}
	r := rep.(*repository)
	defer func() {
		// for test debug
		// drop collection after test complete successful
		/*if err := r.tasksC.Drop(ctx); err != nil {
			log.Error(err)
		}*/
	}()
	address := "0x81b7e08f65bdf5648606c89998a9cc8164397647"
	txId := "0x68392adbfd32cce6170eb909ad8c889319840593692df18c9a1b24818a1cfa1d"
	taskAddress := models.Task{
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Type:           models.OneTime,
		BlockchainType: models.Ethereum,
		ListenTo:       models.ListenObject{Type: models.ListenTypeAddress, Value: address},
		Callback:       models.Callback{Url: "who_waits/address", Type: models.Get},
	}
	taskAddressId, err := rep.PutTask(ctx, &taskAddress)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}

	taskTxHash := models.Task{
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Type:           models.OneTime,
		BlockchainType: models.Ethereum,
		ListenTo:       models.ListenObject{Type: models.ListenTypeTxID, Value: txId},
		Callback:       models.Callback{Url: "who_waits/txId", Type: models.Get},
	}
	taskTxId, err := rep.PutTask(ctx, &taskTxHash)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	// check tasks count in db
	count, err := r.tasksC.Count(ctx, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), count)
	// only address suitable
	tasks, err := rep.FindByAddressOrTxId(ctx, models.Ethereum, address, "0x912dda9f2618d6a1876863e26cf1efa9b7603e4c239aadb19556f16a0d7f8508")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, 1, len(tasks), "test 'get task by address' fail")
	if len(tasks) != 1 {
		t.FailNow()
	}
	assert.Equal(t, taskAddressId, tasks[0].Id.Hex(), "test 'get task by address' fail")
	assert.Equal(t, "who_waits/address", tasks[0].Callback.Url, "test 'get task by address' fail")

	// only txID suitable
	tasks2, err := rep.FindByAddressOrTxId(ctx, models.Ethereum, "0xf4bac4964dafa8d02ce1ace0b75b753b1dde2ac5", txId)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, 1, len(tasks2), "test 'get task by txID' fail")
	if len(tasks2) != 1 {
		t.FailNow()
	}
	assert.Equal(t, taskTxId, tasks2[0].Id.Hex(), "test 'get task by txID' fail")
	assert.Equal(t, "who_waits/txId", tasks2[0].Callback.Url, "test 'get task by txID' fail")

	// address and txID suitable
	tasks3, err := rep.FindByAddressOrTxId(ctx, models.Ethereum, address, txId)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, 2, len(tasks3), "test 'get tasks by txID and address' fail")

	// blockchain is not suitable
	tasks4, err := rep.FindByAddressOrTxId(ctx, models.ChainType("WAVES"), address, txId)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, 0, len(tasks4), "test 'not suitable blockchain' fail")

	// remove task with address
	if err := rep.RemoveTask(ctx, taskAddressId); err != nil {
		log.Error(err)
		t.FailNow()
	}

	// get removed task
	tasks5, err := rep.FindByAddressOrTxId(ctx, models.Ethereum, address, "0x912dda9f2618d6a1876863e26cf1efa9b7603e4c239aadb19556f16a0d7f8508")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, 0, len(tasks5), "test 'remove task' fail")

	// remove task with txId - clean DB after test
	if err := rep.RemoveTask(ctx, taskTxId); err != nil {
		log.Error(err)
		t.FailNow()
	}
}

func beforeTest() (ctx context.Context, log logger.ILogger) {
	ctx = context.Background()
	log, _ = logger.Init(false, logger.DEBUG)
	return
}
