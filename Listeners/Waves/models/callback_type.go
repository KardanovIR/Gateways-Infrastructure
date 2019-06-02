package models

type CallbackType string

const (
	InitInTx      CallbackType = "InitInTx"
	InitOutTx     CallbackType = "InitOutTx"
	FinishProcess CallbackType = "FinishProcess"
	CompleteTx    CallbackType = "CompleteTx"
)
