package services

import (
	"bytes"
	"context"
	"golang.org/x/crypto/blake2b"

	"github.com/btcsuite/btcutil/base58"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
)

func (cl *nodeClient) ValidateAddress(ctx context.Context, address string) (bool, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'ValidateAddress' for %s", address)

	//todo сделать метод
	return true, nil
}

func (cl *nodeClient) PublicKeyFromAddress(ctx context.Context, address string) []byte {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'PublicKeyFromAddress' for %s", address)
	//todo сделать метод
	return  nil
}
