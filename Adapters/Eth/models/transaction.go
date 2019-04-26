package models

type TxStatus string

const (
	TxStatusUnKnown TxStatus = "UNKNOWN"
	TxStatusPending TxStatus = "PENDING"
	TxStatusSuccess TxStatus = "SUCCESS"
)

type TxInfo struct {
	From     string
	To       string
	Amount   string
	Contract string
	TxHash   string
	Fee      string
	Data     string
	Status   TxStatus
}
