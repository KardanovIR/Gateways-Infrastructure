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
	Inputs          []InputOutputInfo
	Outputs         []InputOutputInfo
}

type InputOutputInfo struct {
	Address string
	Amount  string
}

type UnSignedTx struct {
	ID         string        `json:"id,omitempty"`
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

type Tx struct {
	Summary    Summary       `json:"summary"`
	Inputs     []InputOutput `json:"inputs"`
	DataInputs []interface{} `json:"dataInputs"`
	Outputs    []InputOutput `json:"outputs"`
}

type Summary struct {
	ID string `json:"id,omitempty"`
}

type InputOutput struct {
	Address  string `json:"address"`
	Value    uint64 `json:"value"`
	ErgoTree string `json:"ergoTree"`
}
