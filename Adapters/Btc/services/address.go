package services

import (
	"context"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
)

func (cl *nodeClient) ValidateAddress(ctx context.Context, address string) (bool, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'ValidateAddress' for %s", address)
	defaultNet := &chaincfg.MainNetParams
	btcAddress, err := btcutil.DecodeAddress(address,defaultNet)
	if err != nil {
		log.Error(err)
		return false, err
	}

	result, err := cl.nodeClient.ValidateAddress(btcAddress)
	if err != nil {
		log.Error(err)
		return false, err
	}

	return result.IsValid, nil
}