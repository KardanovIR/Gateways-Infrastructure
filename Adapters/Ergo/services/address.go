package services

import (
	"context"

	"github.com/btcsuite/btcutil/base58"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
)

const (
	checksumBytesCount    = 4
	addressTypeBytesCount = 1
)

func (cl *nodeClient) ValidateAddress(ctx context.Context, address string) (bool, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'ValidateAddress' for %s", address)
	// todo implementation
	return true, nil
}

func (cl *nodeClient) PublicKeyFromAddress(ctx context.Context, address string) []byte {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'PublicKeyFromAddress' for %s", address)
	addr := base58.Decode(address)
	publicKey := addr[addressTypeBytesCount : len(addr)-checksumBytesCount]
	return publicKey
}
