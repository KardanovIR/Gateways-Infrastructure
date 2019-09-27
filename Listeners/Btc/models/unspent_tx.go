package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UnspentTx struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Address string             `bson:"address"`
	TxHash  string             `bson:"txHash"`
	Amount  uint64             `bson:"amount"`
	OutputN uint32             `bson:"outputN"`
	Locked  bool               `bson:"locked"`
}
