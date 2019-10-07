package models

type RawUtxo struct {
	TxId            string  `json:"txid,omitempty"`
	Address         string  `json:"address"`
	Vout            uint    `json:"vout"`
	ScriptPubKey    string  `json:"scriptPubKey"`
	Amount          float32 `json:"amount"`
	Satoshis        uint    `json:"satoshis"`
	Height          uint    `json:"height"`
	Confirmations   uint    `json:"confirmations"`
}
