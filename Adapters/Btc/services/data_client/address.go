package data_client

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
	getUnspentInputsUrlTemplate = "/addr/%s/utxo"
)

func (dcl *dataClient) GetAllBalances(ctx context.Context, address string) (*models.Balance, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'GetBalanceWithAssets' for address %s", address)
	balance := models.Balance{}

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


func (dcl *dataClient) GetUnspentInputs(ctx context.Context, address string) ([]*models.RawUtxo, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'GetUnspendInputs' for address %s", address)
	utxo := make([]*models.RawUtxo, 0)

	r, err := dcl.Request(ctx, http.MethodGet, dcl.conf.Url+fmt.Sprintf(getUnspentInputsUrlTemplate, address), nil)
	if err != nil {
		log.Errorf("failed to get balance for address %s: %s", address, err)
		return nil, err
	}
	if err := json.Unmarshal(r, &utxo); err != nil {
		log.Errorf("failed to unmarshal balance: %s", err)
		return nil, err
	}
	return utxo, nil
}

