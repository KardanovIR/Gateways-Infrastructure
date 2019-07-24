package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/services/converter"
)

const (
	getBalanceUrlTemplate = "/addresses/%s"
)

type BalanceResponse struct {
	Transactions AddressTransactions `json:"transactions"`
}

type AddressTransactions struct {
	ConfirmedBalance uint64 `json:"confirmedBalance"`
}

func (cl *nodeClient) GetAllBalances(ctx context.Context, address string) (*models.AccountBalance, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'GetBalanceWithAssets' for address %s", address)
	balance := models.AccountBalance{}

	r, err := cl.Request(ctx, http.MethodGet, cl.conf.ExplorerUrl+fmt.Sprintf(getBalanceUrlTemplate, address), nil)
	if err != nil {
		log.Errorf("failed to get balance for address %s: %s", address, err)
		return nil, err
	}
	balanceResponse := BalanceResponse{}
	if err := json.Unmarshal(r, &balanceResponse); err != nil {
		log.Errorf("failed to unmarshal balance: %s", err)
		return nil, err
	}
	balance.Amount = converter.ToTargetAmount(balanceResponse.Transactions.ConfirmedBalance)
	return &balance, nil
}
