package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/models"
)

const (
	blockAtTemplateUrl = "/blocks/at/%d"
	blockLastUrl       = "/blocks/last"
)

func (service *nodeReader) blockAt(ctx context.Context, height uint64) (*models.Block, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'blockAt' %d", height)
	return service.getBlock(ctx, fmt.Sprintf(blockAtTemplateUrl, height))
}

func (service *nodeReader) blockLast(ctx context.Context) (*models.Block, error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'blockLast'")
	return service.getBlock(ctx, blockLastUrl)
}

func (service *nodeReader) getBlock(ctx context.Context, url string) (*models.Block, error) {
	log := logger.FromContext(ctx)
	urlFull := service.nodeClient.GetOptions().BaseUrl + url
	block := new(models.Block)
	req, err := http.NewRequest("GET", urlFull, nil)
	if err != nil {
		return nil, err
	}
	if _, err := service.nodeClient.Do(ctx, req, block); err != nil {
		log.Errorf("get block fails: %s", err)
		return nil, err
	}
	return block, nil
}
