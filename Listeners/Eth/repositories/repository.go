package repositories

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/models"
)

const (
	Ctasks      = "tasks"
	CChainState = "chainStates"
)

type IRepository interface {
	PutTask(ctx context.Context, task models.Task) (string, error)
	RemoveTask(ctx context.Context, id string) error
	FindByAddress(ctx context.Context, ticket models.ChainType, addresses string) (tasks []models.Task, err error)
	GetLastChainState(ctx context.Context, chainType models.ChainType) (chainState *models.ChainState, err error)
	PutChainState(ctx context.Context, state models.ChainState) (err error)
}

type repository struct {
	session     *mgo.Session
	tasksC      *mgo.Collection
	chainStateC *mgo.Collection
	dbName      string
}

var (
	rep            IRepository
	onceRepository sync.Once
)

func GetRepository() IRepository {
	onceRepository.Do(func() {
		panic("try to get node reader before it's creation!")
	})
	return rep
}

func New(ctx context.Context, url, dbName string) error {
	log := logger.FromContext(ctx)
	var initErr error
	onceRepository.Do(func() {
		for i := 1; i < 6; i++ {
			log.Debug("Attempt %d to connect to MongoDB at %s", i, url)
			conn, err := mgo.Dial(url)
			if err != nil {
				log.Errorf("Failed to connect to MongoDB at %s: %v", url, err)
				initErr = err
				time.Sleep(1 * time.Second)
			} else {
				var db = conn.DB(dbName)
				log.Infof("Connected successfully to MongoDB at %s", url)
				rep = &repository{
					session:     conn,
					tasksC:      db.C(Ctasks),
					chainStateC: db.C(CChainState),
					dbName:      dbName}
			}
		}
	})
	return initErr
}

func (rep *repository) PutTask(ctx context.Context, task models.Task) (id string, err error) {
	log := logger.FromContext(ctx)
	log.Info("InsertTask")

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

	if string(id) == "" {
		var objectId = bson.ObjectId(id)
		err = rep.tasksC.Remove(bson.M{
			"_id": objectId,
		})
		if err != nil {
			log.Errorf("Inserting task to DB fails: %s", err)
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

func (rep *repository) GetLastChainState(ctx context.Context, chainType models.ChainType) (chainState *models.ChainState, err error) {
	log := logger.FromContext(ctx)
	log.Info("GetLastChainState")
	rep.refreshSession()

	err = rep.chainStateC.Find(bson.M{"timestamp": -1, "chaintype": chainType}).One(&chainState)
	if err != nil {
		log.Errorf("Getting lastChainState from DB fails: %s", err)
	}
	return
}

func (rep *repository) PutChainState(ctx context.Context, state models.ChainState) (err error) {
	log := logger.FromContext(ctx)
	log.Info("PutChainState")
	rep.refreshSession()

	if string(state.Id) == "" {
		state.Id = bson.NewObjectId()
		err = rep.chainStateC.Insert(state)
		return
	}

	_, err = rep.chainStateC.UpsertId(state.Id, state)

	return

}
func (rep *repository) refreshSession() {
	rep.session.Refresh()
	// todo: maybe this is overhead
	db := rep.session.DB(rep.dbName)
	rep.tasksC = db.C(Ctasks)
}
