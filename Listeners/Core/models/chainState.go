package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type ChainState struct {
	Id        bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
	ChainType ChainType     `json:"chaintype" bson:"chaintype"`
	LastBlock int64         `json:"lastblock" bson:"lastblock"`
}
