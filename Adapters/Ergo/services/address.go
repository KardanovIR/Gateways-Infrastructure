package services

import (
	"context"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
)

func (cl *nodeClient) ValidateAddress(ctx context.Context, address string) (bool, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'ValidateAddress' for %s", address)
	// todo implementation
	return true, nil
}
