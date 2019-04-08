package service

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strings"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Router/grpc/blockchain"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/model"
)

const (
	headerName = "blockchain-type"
)

func (bs *blockchainsService) GetTransactionStatus(ctx context.Context, blockchain model.Blockchain, txID string) (string, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'GetTransactionStatus' for %s, txID: %s", blockchain, txID)
	requestCtx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(headerName,
		strings.ToLower(string(blockchain))))
	request := pb.GetTransactionStatusRequest{TxId: txID}
	reply, err := bs.universal.GetTransactionStatus(requestCtx, &request)
	if err != nil {
		return "", err
	}
	return reply.Status, nil
}
