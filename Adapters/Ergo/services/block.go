package services

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
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
	r, _ := cl.Request(ctx, http.MethodGet, cl.conf.ExplorerUrl+getCurrentBlockUrl, nil)
	getCurrentBlockResp := BlocksResponse{}
	if err := json.Unmarshal(r, &getCurrentBlockResp); err != nil {
		log.Errorf("failed to get current height: %s", err)
		return 0, err
	}
	height := getCurrentBlockResp.Items[0].Height
	log.Infof("current height is %d", height)
	return height, nil
}
