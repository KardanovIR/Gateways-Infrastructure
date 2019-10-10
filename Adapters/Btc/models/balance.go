package models

type Balance struct {
	Address string `bson:"address"`
	Amount  uint64 `bson:"amount"`
}
