package repositories

import (
	"context"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/config"
	"sync"
	"time"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IRepository interface {
	GetUnspentTxListForAddresses(ctx context.Context, addresses []string) ([]models.UnspentTx, error)
	GetUnspentTxByTxHashAndOutputNumber(ctx context.Context, txHash string, outputN uint32) (*models.UnspentTx, error)
	GetBalanceForAddresses(ctx context.Context, addresses []string) ([]models.Balance, error)
}

const UnspentTxCollection = "unspentTx"

type repository struct {
	client     *mongo.Client
	unspentTxC *mongo.Collection
	DbName     string
}

// if addresses is empty return all inputs
func (r *repository) GetUnspentTxListForAddresses(ctx context.Context, addresses []string) ([]models.UnspentTx, error) {
	log := logger.FromContext(ctx)
	log.Debugf("GetUnspentTxListForAddresses %s", addresses)
	filter := bson.D{{}}
	if len(addresses) > 0 {
		filter = bson.D{{"address", bson.D{{"$in", addresses}}}}
	}
	cur, err := r.unspentTxC.Find(ctx, filter)
	if err != nil {
		log.Errorf("Finding unspent tx for %s in DB fails: %s", addresses, err)
		return nil, err
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			log.Error("close cursor error: ", err)
		}
	}()
	unspentTxList := make([]models.UnspentTx, 0)
	for cur.Next(ctx) {
		var unspent models.UnspentTx
		if err := cur.Decode(&unspent); err != nil {
			return unspentTxList, err
		}
		log.Debugf("Unspent tx %+v", unspent)
		unspentTxList = append(unspentTxList, unspent)
	}
	if err := cur.Err(); err != nil {
		return unspentTxList, err
	}
	return unspentTxList, nil
}

func (r *repository) GetBalanceForAddresses(ctx context.Context, addresses []string) ([]models.Balance, error) {
	log := logger.FromContext(ctx)
	log.Debugf("GetBalanceForAddresses %s", addresses)
	stageMatch := bson.D{{
		"$match", bson.D{{
			"address", bson.D{{"$in", addresses}},
		}},
	}}
	stageSumByAddress := bson.D{{"$group", bson.D{{"_id", "$address"}, {"amount", bson.D{{"$sum", "$amount"}}}}}}
	stageRenameFields := bson.D{{"$project", bson.M{"address": "$_id", "amount": 1}}}
	var pipeline bson.A
	if len(addresses) == 0 {
		// collect by all addresses
		pipeline = bson.A{stageSumByAddress, stageRenameFields}
	} else {
		pipeline = bson.A{stageMatch, stageSumByAddress, stageRenameFields}
	}
	cur, err := r.unspentTxC.Aggregate(ctx, pipeline)

	if err != nil {
		log.Errorf("Finding unspent tx for %s in DB fails: %s", addresses, err)
		return nil, err
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			log.Error("close cursor error: ", err)
		}
	}()
	balances := make([]models.Balance, 0)

	for cur.Next(ctx) {
		var balance models.Balance
		if err := cur.Decode(&balance); err != nil {
			return balances, err
		}
		balances = append(balances, balance)
	}
	if err := cur.Err(); err != nil {
		return balances, err
	}
	return balances, nil
}

func (r *repository) GetUnspentTxByTxHashAndOutputNumber(ctx context.Context, txHash string, outputN uint32) (*models.UnspentTx, error) {
	log := logger.FromContext(ctx)
	log.Debugf("GetUnspentTxByTxHashAndOutputNumber %s, %d", txHash, outputN)
	result := r.unspentTxC.FindOne(ctx, bson.D{{
		"$and", bson.A{
			bson.D{{"txHash", txHash}},
			bson.D{{"outputN", outputN}},
		}},
	})
	var unspentTx models.UnspentTx
	err := result.Decode(&unspentTx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Errorf("GetUnspentTxByTxHashAndOutputNumber from DB fails: %s", err)
		return nil, err
	}
	return &unspentTx, nil
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

func New(ctx context.Context, dbConf config.DB) error {
	var initErr error
	onceRepository.Do(func() {
		rep, initErr = connect(ctx, dbConf.Host, dbConf.Name)
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
			rep = &repository{
				client:     mongoClient,
				unspentTxC: db.Collection(UnspentTxCollection),
				DbName:     dbName}
			return rep, nil
		}
	}
	return nil, err
}
