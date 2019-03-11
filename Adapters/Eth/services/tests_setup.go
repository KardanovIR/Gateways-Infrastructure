// utilities for tests
package services

import (
	"github.com/ethereum/go-ethereum/crypto"
)

func SetKeyPair(address string, pkHex string) error {
	pk, err := crypto.HexToECDSA(pkHex)
	if err != nil {
		return err
	}
	cl.(*nodeClient).privateKeys[address] = pk
	return nil
}
