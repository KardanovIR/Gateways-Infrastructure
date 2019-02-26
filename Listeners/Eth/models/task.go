package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Task struct {
	Id              bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	CreatedAt       time.Time       `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt" bson:"updatedAt"`
	Address 		string          `json:"address" bson:"address"`
	CallbackUrl     string          `json:"callbackUrl" bson:"callbackUrl"`
	Type            TaskType        `json:"taskType" bson:"taskType"`
	BlockchainType	BlockchainType	`json:"blockchainType" bson:"blockchainType"`
}

type TaskType int

const (
	Permanent TaskType = iota + 1
	OneTime
)

type BlockchainType string

const (
	Ethereum BlockchainType = "Ethereum"
)
