package models

import "math/big"

// balance for account: eth and tokens
type AccountBalance struct {
	// eth amount
	Amount *big.Int
	// contract's address of token - token's balance
	Tokens map[string]*big.Int
}
