package repository

import (
	"context"
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUTXORepository interface {
	repositories.IRepository
	DeleteUnspentTx(ctx context.Context, txHash string, outputN uint32) (*models.UnspentTx, error)
	AddUnspentTx(ctx context.Context, unspentTx models.UnspentTx) error
	GetUnspentTxListForAddress(ctx context.Context, address string) ([]models.UnspentTx, error)
	GetUnspentTxByTxHashAndOutputNumber(ctx context.Context, txHash string, outputN uint32) (*models.UnspentTx, error)
}

const UnspentTxCollection = "unspentTx"

type repository struct {
	*repositories.Repository
	unspentTxC *mongo.Collection
}

func (r *repository) DeleteUnspentTx(ctx context.Context, txHash string, outputN uint32) (*models.UnspentTx, error) {
	log := logger.FromContext(ctx)
	log.Debugf("DeleteUnspentTx %s N = %d", txHash, outputN)
	result := r.unspentTxC.FindOneAndDelete(ctx, bson.D{{
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

func (r *repository) AddUnspentTx(ctx context.Context, unspentTx models.UnspentTx) error {
	log := logger.FromContext(ctx)
	log.Debugf("Add input %+v", unspentTx)
	if unspentTx.Id.IsZero() {
		_, err := r.unspentTxC.InsertOne(ctx, &unspentTx)
		if err != nil {
			log.Errorf("Inserting unspentTx to DB fails: %s", err)
			return err
		}
		return nil
	} else {
		return fmt.Errorf("unspent input %s is already exist", unspentTx.Id.Hex())
	}
}

// todo don't used -> move to adapter
func (r *repository) GetUnspentTxListForAddress(ctx context.Context, address string) ([]models.UnspentTx, error) {
	log := logger.FromContext(ctx)
	log.Debugf("GetUnspentTxListForAddress %s", address)
	cur, err := r.unspentTxC.Find(ctx, bson.D{{"address", address}})
	if err != nil {
		log.Errorf("Finding unspent tx for %s in DB fails: %s", address, err)
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

func New(ctx context.Context, dbConfig config.DB) (IUTXORepository, error) {
	log := logger.FromContext(ctx)
	err := repositories.New(ctx, dbConfig.Host, dbConfig.Name)
	if err != nil {
		log.Error("Can't create db connection: ", err)
		return nil, err
	}
	rb := repositories.GetRepository().(*repositories.Repository)
	unspentTxC := rb.Client.Database(rb.DbName).Collection(UnspentTxCollection)
	rep := &repository{rb, unspentTxC}
	return rep, nil
}
