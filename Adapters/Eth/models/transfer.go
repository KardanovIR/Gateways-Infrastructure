package models

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// TransferEvent is used by abi.unpack
type TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

func (e *TransferEvent) String() string {
	return fmt.Sprintf("from, to, value: %s, %s, %s",
		e.From.String(), e.To.String(), e.Value.String())
}
