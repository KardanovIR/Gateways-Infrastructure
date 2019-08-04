package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/objectid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/options"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
)

func (rep *repository) GetLastChainState(ctx context.Context, chainType models.ChainType) (*models.ChainState, error) {
	log := logger.FromContext(ctx)
	log.Debugf("GetLastChainState for %s", chainType)
	var chainState = new(models.ChainState)
	result := rep.chainStateC.FindOne(ctx, bson.D{{"chaintype", chainType}})
	err := result.Decode(chainState)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Errorf("Getting lastChainState from DB fails: %s", err)
	}
	return chainState, err
}

func (rep *repository) PutChainState(ctx context.Context, state *models.ChainState) (newState *models.ChainState, err error) {
	log := logger.FromContext(ctx)
	log.Debugf("PutChainState %d", state.LastBlock)
	state.Timestamp = time.Now()
	insert := true // if don't find chain state - insert it
	ur, err := rep.chainStateC.ReplaceOne(ctx, bson.D{{"chaintype", state.ChainType}}, state, &options.ReplaceOptions{Upsert: &insert})
	if err != nil {
		return nil, err
	}
	if ur.UpsertedID != nil {
		state.Id = ur.UpsertedID.(objectid.ObjectID)
	}
	return state, err
}
