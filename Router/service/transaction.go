package service

import (
	"context"
	"fmt"

	"github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/ethAdapter"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/wavesAdapter"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/model"
)

func (bs *blockchainsService) GetTransactionStatus(ctx context.Context, blockchain model.Blockchain, txID string) (string, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'GetTransactionStatus' for %s, txID: %s", blockchain, txID)
	switch blockchain {
	case model.ETH:
		request := ethAdapter.GetTransactionStatusRequest{TxHash: txID}
		reply, err := bs.ethAdapter.GetTransactionStatus(ctx, &request)
		if err != nil {
			return "", err
		}
		return reply.Status, nil
	case model.WAVES:
		request := wavesAdapter.GetTransactionStatusRequest{TxId: txID}
		reply, err := bs.wavesAdapter.GetTransactionStatus(ctx, &request)
		if err != nil {
			return "", err
		}
		return reply.Status, nil
	}
	return "", fmt.Errorf("not supported blockchain %s", blockchain)
}
