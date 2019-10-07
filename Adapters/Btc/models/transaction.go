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

type RawTx struct {
	Id            string      `json:"txid,omitempty"`
	Inputs        []RawInput  `json:"vin"`
	Outputs       []RawOutput `json:"vout"`
	Version       uint        `json:"version"`
	BlockHash     string      `json:"blockhash"`
	BlockHeight   uint        `json:"blockheight"`
	Confirmations uint        `json:"confirmations"`
	Time          uint        `json:"time"`
	BlockTime     uint        `json:"blocktime"`
	ValueOut      float64     `json:"valueOut"`
	ValueIn       float64     `json:"valueIn"`
	Fees          float64     `json:"fees"`
	Size          uint        `json:"size"`
	LockTime      uint        `json:"locktime"`
}

type RawInput struct {
	Id              string  `json:"txid,omitempty"`
	Vout            uint    `json:"vout"`
	Sequence        uint    `json:"sequence"`
	N               uint    `json:"n"`
	Script          Script  `json:"scriptSig"`
	Address         string  `json:"addr"`
	ValueSat        uint    `json:"valueSat"`
	Value           float64 `json:"value"`
	DoubleSpentTxID string  `json:"doubleSpentTxID"`
}

type RawOutput struct {
	Value       string `json:"value"`
	N           uint   `json:"n"`
	Script      Script `json:"scriptPubKey"`
	SpentTxId   string `json:"spentTxId"`
	SpentIndex  string `json:"spentIndex"`
	SpentHeight uint   `json:"spentHeight"`
}

type Script struct {
	Asm        string   `json:"asm"`
	Hex        string   `json:"hex"`
	ScriptType string   `json:"type"`
	Addresses  []string `json:"addresses"`
}

