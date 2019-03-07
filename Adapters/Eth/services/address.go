package services

import (
	"context"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
)

func (cl *nodeClient) GenerateAddress(ctx context.Context) (publicAddress string, err error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'GenerateAddress'")
	prkey, err := crypto.GenerateKey()
	if err != nil {
		log.Errorf("Failed to create new private key %v", err)
		return "", err
	}
	prKeyHex := hex.EncodeToString(crypto.FromECDSA(prkey))
	addr := crypto.PubkeyToAddress(prkey.PublicKey)
	publicAddress = addr.Hex()
	log.Infof("Private hex %s, public address %s", prKeyHex, publicAddress)
	cl.privateKeys[publicAddress] = prkey
	return
}

func (cl *nodeClient) IsAddressValid(ctx context.Context, address string) bool {
	if !common.IsHexAddress(address) {
		return false
	}
	return true
}
