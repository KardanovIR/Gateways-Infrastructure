package services

import (
	"context"
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
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

	//todo сделать метод

	return nil, fmt.Errorf("not inplemented")
}
