package node

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

func (cl *nodeClient) getCurrentHeight(ctx context.Context) (int64, error) {
	log := logger.FromContext(ctx)
	log.Info("get current height")
	height, err := cl.nodeClient.GetBlockCount()
	if err != nil {
		return 0, err
	}

	log.Infof("current height is %d", height)
	return height, nil
}

