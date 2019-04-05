package repositories

import (
	"context"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/models"
)

const (
	Ctasks      = "tasks"
	CChainState = "chainStates"
)

type IRepository interface {
	PutTask(ctx context.Context, task *models.Task) (string, error)
	RemoveTask(ctx context.Context, id string) error
	FindByAddressOrTxId(ctx context.Context, ticket models.ChainType, address string, txID string) (tasks []*models.Task, err error)
	GetLastChainState(ctx context.Context, chainType models.ChainType) (chainState *models.ChainState, err error)
	PutChainState(ctx context.Context, state *models.ChainState) (newState *models.ChainState, err error)
}

type repository struct {
	client      *mongo.Client
	tasksC      *mongo.Collection
	chainStateC *mongo.Collection
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
	var initErr error
	onceRepository.Do(func() {
		rep, initErr = connect(ctx, url, dbName)
	})
	return initErr
}

func connect(ctx context.Context, url, dbName string) (IRepository, error) {
	log := logger.FromContext(ctx)
	mongoClient, err := mongo.Connect(ctx, "mongodb://"+url)
	if err != nil {
		log.Errorf("Failed to create connect configuration to MongoDB %s: %v", url, err)
		return nil, err
	}
	for i := 1; i < 6; i++ {
		log.Debugf("Attempt %d to connect to MongoDB at %s", i, url)
		err := mongoClient.Ping(ctx, nil)
		if err != nil {
			log.Errorf("Failed to connect to MongoDB at %s: %v", url, err)
			time.Sleep(1 * time.Second)
		} else {
			var db = mongoClient.Database(dbName)
			log.Infof("Connected successfully to MongoDB at %s", url)
			rep = &repository{
				client:      mongoClient,
				tasksC:      db.Collection(Ctasks),
				chainStateC: db.Collection(CChainState),
				dbName:      dbName}
			return rep, nil
		}
	}
	return nil, err
}
