package models

// balance for account: waves and assets
type AccountBalance struct {
	// waves amount
	Amount uint64
	// assetId in Base58 - asset's balance
	Assets map[string]uint64
}
