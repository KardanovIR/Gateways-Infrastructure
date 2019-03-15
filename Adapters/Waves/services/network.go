package services

import (
	"context"
	"strconv"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
)

func (cl *nodeClient) GetLastBlockHeight(ctx context.Context) (string, error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'GetLastBlockHeight'")

	lastBlock, _, err := cl.nodeClient.Blocks.Last(ctx)
	if err != nil {
		log.Errorf("get last block fails: %s", err)
		return "", err
	}
	return strconv.FormatUint(lastBlock.Height, 10), nil
}

func (cl *nodeClient) Fee(ctx context.Context) (uint64, error) {
	return 100000, nil
}
