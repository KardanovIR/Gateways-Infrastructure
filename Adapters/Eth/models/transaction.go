package models

type TxStatus string

const (
	TxStatusUnKnown TxStatus = "UNKNOWN"
	TxStatusPending TxStatus = "PENDING"
	TxStatusSuccess TxStatus = "SUCCESS"
)
