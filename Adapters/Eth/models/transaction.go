package models

import "math/big"

type TxStatus string

const (
	TxStatusUnKnown TxStatus = "UNKNOWN"
	TxStatusPending TxStatus = "PENDING"
	TxStatusSuccess TxStatus = "SUCCESS"
)

type TxInfo struct {
	From     string
	To       string
	Amount   *big.Int
	Contract string
	TxHash   string
	Fee      *big.Int
	Data     string
	Status   TxStatus
}
