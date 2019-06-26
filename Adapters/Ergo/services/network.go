package services

import (
	"context"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
)

func (cl *nodeClient) Fee(ctx context.Context, senderPublicKey string, assetId string) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'Fee' for sender %s, assetId %s", senderPublicKey, assetId)
	return txFee, nil
}
