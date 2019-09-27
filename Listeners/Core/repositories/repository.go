package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"

	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
	"go.mongodb.org/mongo-driver/mongo"
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

type Repository struct {
	Client      *mongo.Client
	tasksC      *mongo.Collection
	chainStateC *mongo.Collection
	DbName      string
}

var (
	rep            IRepository
	onceRepository sync.Once
)

func GetRepository() IRepository {
	onceRepository.Do(func() {
		panic("try to get Repository before it's creation!")
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
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+url))
	if err != nil {
		log.Errorf("Failed to create connect configuration to MongoDB %s: %v", url, err)
		return nil, err
	}
	for i := 1; i < 6; i++ {
		log.Infof("Attempt %d to connect to MongoDB at %s", i, url)
		err := mongoClient.Ping(ctx, nil)
		if err != nil {
			log.Errorf("Failed to connect to MongoDB at %s: %v", url, err)
			time.Sleep(1 * time.Second)
		} else {
			var db = mongoClient.Database(dbName)
			log.Infof("Connected successfully to MongoDB at %s", url)
			rep = &Repository{
				Client:      mongoClient,
				tasksC:      db.Collection(Ctasks),
				chainStateC: db.Collection(CChainState),
				DbName:      dbName}
			return rep, nil
		}
	}
	return nil, err
}
