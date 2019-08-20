package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChainState struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	ChainType ChainType          `json:"chaintype" bson:"chaintype"`
	LastBlock int64              `json:"lastblock" bson:"lastblock"`
}
