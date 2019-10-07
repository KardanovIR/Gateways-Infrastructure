package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
	"net/http"
)

const (
	getBalanceUrlTemplate = "/addr/%s/balance"
)

func (dcl *dataClient) GetAllBalances(ctx context.Context, address string) (*models.AccountBalance, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'GetBalanceWithAssets' for address %s", address)
	balance := models.AccountBalance{}

	r, err := dcl.Request(ctx, http.MethodGet, dcl.conf.Url+fmt.Sprintf(getBalanceUrlTemplate, address), nil)
	if err != nil {
		log.Errorf("failed to get balance for address %s: %s", address, err)
		return nil, err
	}
	if err := json.Unmarshal(r, &balance.Amount); err != nil {
		log.Errorf("failed to unmarshal balance: %s", err)
		return nil, err
	}
	return &balance, nil
}
