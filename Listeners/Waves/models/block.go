package models

type Block struct {
	Height       uint64            `json:"height"`
	Transactions TransactionsField `json:"transactions"`
}
