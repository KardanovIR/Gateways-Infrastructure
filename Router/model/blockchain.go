package model

type Blockchain string

const (
	ETH   Blockchain = "ETH"
	WAVES Blockchain = "WAVES"
)

func (b Blockchain) Exist() bool {
	return b == ETH || b == WAVES
}
