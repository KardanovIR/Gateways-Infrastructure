package repositories

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/models"
)

func (rep *repository) PutTask(ctx context.Context, task models.Task) (id string, err error) {
	log := logger.FromContext(ctx)
	log.Infof("InsertTask %+v", task)

	rep.refreshSession()

	if string(task.Id) == "" {
		task.Id = bson.NewObjectId()
		err = rep.tasksC.Insert(task)
		if err != nil {
			log.Errorf("Inserting task to DB fails: %s", err)
			return "", err
		}
		return task.Id.String(), nil
	}

	err = rep.tasksC.UpdateId(task.Id, task)

	return task.Id.String(), nil
}

func (rep *repository) RemoveTask(ctx context.Context, id string) (err error) {
	log := logger.FromContext(ctx)
	log.Info("RemoveTask")
	rep.refreshSession()

	if string(id) != "" {
		var objectId = bson.ObjectId(id)
		err = rep.tasksC.Remove(bson.M{
			"_id": objectId,
		})
		if err != nil {
			log.Errorf("Removing task to DB fails: %s", err)
			return err
		}
		return nil
	}

	return nil
}

func (rep *repository) FindByAddress(ctx context.Context, ticket models.ChainType, addresses string) (tasks []models.Task, err error) {
	log := logger.FromContext(ctx)
	log.Infof("FindByAddress %s", addresses)
	rep.refreshSession()

	err = rep.tasksC.Find(bson.M{
		"address":        addresses,
		"blockchainType": ticket,
	}).All(&tasks)
	if err != nil {
		log.Errorf("Finding task in DB fails: %s", err)
		return nil, err
	}

	return
}
