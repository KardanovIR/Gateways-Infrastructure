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
	Data            string
	Status          TxStatus
}

type UnSignedTx struct {
	Inputs     []TxInput     `json:"inputs"`
	DataInputs []interface{} `json:"dataInputs"`
	Outputs    []TxOutput    `json:"outputs"`
}

type TxInput struct {
	BoxId         string        `json:"boxId"`
	SpendingProof SpendingProof `json:"spendingProof"`
}

type SpendingProof struct {
	ProofBytes string      `json:"proofBytes"`
	Extension  interface{} `json:"extension"`
}

type TxOutput struct {
	ErgoTree            string        `json:"ergoTree"`
	Assets              []interface{} `json:"assets"`
	AdditionalRegisters interface{}   `json:"additionalRegisters"`
	Value               uint64        `json:"value"`
	CreationHeight      uint64        `json:"creationHeight"`
}

type UnSpentTx struct {
	ID                  string        `json:"id"`
	Value               uint64        `json:"value"`
	CreationHeight      uint64        `json:"creationHeight"`
	ErgoTree            string        `json:"ergoTree"`
	Assets              []interface{} `json:"assets"`
	AdditionalRegisters interface{}   `json:"additionalRegisters"`
	Address             string        `json:"address"`
	SpentTransactionID  string        `json:"spentTransactionId"`
	MainChain           bool          `json:"mainChain"`
}
