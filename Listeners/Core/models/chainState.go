package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type ChainState struct {
	Id        objectid.ObjectID `json:"_id" bson:"_id,omitempty"`
	Timestamp time.Time         `json:"timestamp" bson:"timestamp"`
	ChainType ChainType         `json:"chaintype" bson:"chaintype"`
	LastBlock int64             `json:"lastblock" bson:"lastblock"`
}
