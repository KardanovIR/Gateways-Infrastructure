package services

import (
	"context"
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/services/converter"
)

// CreateRawTx creates transaction
func (cl *nodeClient) CreateRawTx(ctx context.Context, addressFrom string, outs []*models.Output) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'CreateRawTx' from %s to %v",
		addressFrom, outs)
	if len(addressFrom) == 0 || len(outs) == 0 {
		return nil, fmt.Errorf("wrong parameters addressFrom %s or outs", addressFrom)
	}
	amount := uint64(0)
	for _, r := range outs {
		r.Amount = converter.ToNodeAmount(r.Amount)
		amount += r.Amount
	}
	//todo сделать метод
	return nil, nil
}
