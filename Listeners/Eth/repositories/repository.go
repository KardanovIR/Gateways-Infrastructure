package repositories

import (
	"context"
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
	PutChainState(ctx context.Context, state models.ChainState) (newState models.ChainState, err error)
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
		panic("try to get repository before it's creation!")
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

func (rep *repository) refreshSession() {
	rep.session.Refresh()
	// todo: maybe this is overhead
	db := rep.session.DB(rep.dbName)
	rep.tasksC = db.C(Ctasks)
}
