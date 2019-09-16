package services

import (
	"context"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
)

type BlocksResponse struct {
	Items []Block `json:"items"`
}

type Block struct {
	Height uint64 `json:"height"`
}

func (cl *nodeClient) getCurrentHeight(ctx context.Context) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Info("get current height")
	info, err := cl.nodeClient.GetInfo()
	if err != nil {
		return 0, err
	}
	height := uint64(info.Blocks)

	log.Infof("current height is %d", height)
	return height, nil
}
