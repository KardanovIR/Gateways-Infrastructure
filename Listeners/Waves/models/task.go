package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type Task struct {
	Id             objectid.ObjectID `bson:"_id,omitempty"`
	CreatedAt      time.Time         `bson:"createdAt"`
	UpdatedAt      time.Time         `bson:"updatedAt"`
	ListenTo       ListenObject      `bson:"listenTo"`
	Callback       Callback          `bson:"callback"`
	Type           TaskType          `bson:"taskType"`
	BlockchainType ChainType         `bson:"blockchainType"`
}

type Callback struct {
	Type      CallbackType `bson:"callbackType"`
	ProcessId string       `bson:"processId"`
}

// what is listen to: txId or address and it's value
type ListenObject struct {
	Type  ListenType `bson:"type"`
	Value string     `bson:"value"`
}

type ListenType string

const (
	ListenTypeAddress ListenType = "Address"
	ListenTypeTxID    ListenType = "TxId"
)
