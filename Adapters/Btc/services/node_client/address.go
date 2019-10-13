package node_client

import (
	"context"
	"github.com/btcsuite/btcutil"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
)

func (cl *nodeClient) ValidateAddress(ctx context.Context, address string) (bool, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'ValidateAddress' for %s", address)
	btcAddress, err := btcutil.DecodeAddress(address, cl.conf.ChainParams)
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

