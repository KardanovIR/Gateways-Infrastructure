package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Task struct {
	Id             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt" bson:"updatedAt"`
	Address        string        `json:"address" bson:"address"`
	Callback       Callback      `json:"callback" bson:"callback"`
	Type           TaskType      `json:"taskType" bson:"taskType"`
	BlockchainType ChainType     `json:"blockchainType" bson:"blockchainType"`
}

type Callback struct {
	Url  string                 `json:"url" bson:"url"`
	Type CallbackType           `json:"callbackType" bson:"callbackType"`
	Data map[string]interface{} `json:"data" bson:"-"`
}
