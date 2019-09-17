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
	log.Debugf("Private hex %s, public address %s", prKeyHex, publicAddress)
	cl.privateKeys[publicAddress] = prkey
	return
}

func (cl *nodeClient) AddressByPublicKey(ctx context.Context, public string) (string, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'AddressByPublicKey' %s", public)

	publicBytes, err := hex.DecodeString(public)
	if err != nil {
		log.Errorf("Failed convert public key to bytes %s: %s", public, err)
		return "", err
	}
	publicKey, err := crypto.UnmarshalPubkey(publicBytes)
	if err != nil {
		log.Errorf("failed to unmarshal public key %s: %s", public, err)
		return "", err
	}
	addr := crypto.PubkeyToAddress(*publicKey)
	return addr.Hex(), nil
}

// check recipient address
func (cl *nodeClient) IsAddressValid(ctx context.Context, address string) (bool, string, error) {
	if !common.IsHexAddress(address) {
		return false, "", nil
	}
	addr := common.HexToAddress(address)
	isContract, err := cl.IsContract(ctx, addr)
	if err != nil {
		return false, "", err
	}
	// is it smart contract
	if isContract {
		return false, "address is smart contract", err
	}
	return true, "", nil
}

// check address
func (cl *nodeClient) IsContract(ctx context.Context, address common.Address) (bool, error) {
	code, err := cl.ethClient.CodeAt(ctx, address, nil)
	if err != nil {
		return false, err
	}
	// is it smart contract
	if len(code) > 0 {
		return true, nil
	}
	return false, nil
}
