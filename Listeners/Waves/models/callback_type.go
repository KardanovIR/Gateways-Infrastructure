package models

type CallbackType string

const (
	StartProcess  CallbackType = "StartProcess"
	InitInTx      CallbackType = "InitInTx"
	InitOutTx     CallbackType = "InitOutTx"
	FinishProcess CallbackType = "FinishProcess"
)
