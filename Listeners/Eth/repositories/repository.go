package repositories

import (
	"context"
	"time"

	"github.com/globalsign/mgo"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/logger"
)

type IRepository interface {
}

type repository struct {
	dbConnect *mgo.Database
}

func New(ctx context.Context, url, dbName string) (IRepository, error) {
	log := logger.FromContext(ctx)
	var err error
	var conn *mgo.Session
	for i := 1; i < 6; i++ {
		log.Debug("Attempt %d to connect to MongoDB at %s", i, url)
		conn, err = mgo.Dial(url)
		if err != nil {
			log.Errorf("Failed to connect to MongoDB at %s: %v", url, err)
			time.Sleep(1 * time.Second)
		} else {
			log.Infof("Connected successfully to MongoDB at %s", url)
			return &repository{dbConnect: conn.DB(dbName)}, nil
		}
	}
	return nil, err
}
