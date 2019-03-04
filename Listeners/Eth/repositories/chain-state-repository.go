package repositories

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/models"
)

func (rep *repository) GetLastChainState(ctx context.Context, chainType models.ChainType) (chainState *models.ChainState, err error) {
	log := logger.FromContext(ctx)
	log.Info("GetLastChainState")
	rep.refreshSession()

	err = rep.chainStateC.Find(bson.M{"chaintype": chainType}).Sort("timestamp").One(&chainState)
	if err != nil {
		log.Errorf("Getting lastChainState from DB fails: %s", err)
	}
	return
}

func (rep *repository) PutChainState(ctx context.Context, state models.ChainState) (newState models.ChainState, err error) {
	log := logger.FromContext(ctx)
	log.Info("PutChainState")
	rep.refreshSession()

	if string(state.Id) == "" {
		state.Id = bson.NewObjectId()
		err = rep.chainStateC.Insert(state)
		return state, err
	}

	_, err = rep.chainStateC.UpsertId(state.Id, state)

	return state, err

}
