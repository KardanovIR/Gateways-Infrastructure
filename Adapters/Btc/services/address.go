package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
	"net/http"
)

const (
	getBalanceUrlTemplate = "/addr/%s/balance"
	getUnspentInputsUrlTemplate = "/addr/%s/utxo"
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
