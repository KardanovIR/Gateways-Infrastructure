package models

type TxStatus string

const (
	TxStatusUnKnown TxStatus = "UNKNOWN"
	TxStatusPending TxStatus = "PENDING"
	TxStatusSuccess TxStatus = "SUCCESS"
)

type TxInfo struct {
	SenderPublicKey string
	From            string
	To              string
	Amount          string
	AssetId         string
	TxHash          string
	Fee             string
	Data            []byte
	Status          TxStatus
}
