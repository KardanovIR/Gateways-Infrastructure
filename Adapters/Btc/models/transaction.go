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
	Inputs          []InputOutputInfo
	Outputs         []InputOutputInfo
}

type InputOutputInfo struct {
	Address string
	Amount  string
}

type InputOutput struct {
	Address string `json:"address"`
	Value   uint64 `json:"value"`
}
