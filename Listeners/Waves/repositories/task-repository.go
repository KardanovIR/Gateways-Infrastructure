package repositories

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/models"
)

func (rep *repository) PutTask(ctx context.Context, task *models.Task) (id string, err error) {
	log := logger.FromContext(ctx)
	log.Infof("PutTask %+v", task)
	if task.Id.IsZero() {
		ir, err := rep.tasksC.InsertOne(ctx, task)
		if err != nil {
			log.Errorf("Inserting task to DB fails: %s", err)
			return "", err
		}
		id := ir.InsertedID.(objectid.ObjectID)
		task.Id = id
		return id.Hex(), nil
	}
	return "", fmt.Errorf("task id is not empty: %s", id)
}

func (rep *repository) RemoveTask(ctx context.Context, id string) (err error) {
	log := logger.FromContext(ctx)
	log.Infof("RemoveTask %s", id)
	if id != "" {
		objId, err := objectid.FromHex(id)
		if err != nil {
			log.Errorf("task id %s has wrong format: %s", id, err)
			return err
		}
		if _, err := rep.tasksC.DeleteOne(ctx, bson.M{"_id": objId}); err != nil {
			log.Errorf("Removing task from DB fails: %s", err)
			return err
		}
	}
	return nil
}

func (rep *repository) FindByAddressOrTxId(ctx context.Context, ticket models.ChainType, address string,
	txID string) ([]*models.Task, error) {

	log := logger.FromContext(ctx)
	log.Infof("FindByAddressOrTxId: address %s, txID %s", address, txID)
	cur, err := rep.tasksC.Find(ctx,
		bson.D{{
			"$or", bson.A{
				bson.D{{"listenTo", bson.D{{"type", models.ListenTypeAddress}, {"value", address}}}},
				bson.D{{"listenTo", bson.D{{"type", models.ListenTypeTxID}, {"value", txID}}}},
			}},
			{"blockchainType", ticket},
		})
	if err != nil {
		log.Errorf("Finding task in DB fails: %s", err)
		return nil, err
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			log.Error("close cursor error: ", err)
		}
	}()
	tasks := make([]*models.Task, 0)
	for cur.Next(ctx) {
		var task models.Task
		if err := cur.Decode(&task); err != nil {
			return tasks, err
		}
		log.Debugf("%+v", task)
		tasks = append(tasks, &task)
	}
	if err := cur.Err(); err != nil {
		return tasks, err
	}
	return tasks, nil
}
