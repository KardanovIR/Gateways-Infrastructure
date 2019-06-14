package services

import (
	"context"
	"sync"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/models"
)

type INodeClient interface {
	ValidateAddress(ctx context.Context, address string) (bool, error)
	GetAllBalances(ctx context.Context, address string) (*models.AccountBalance, error)

	Fee(ctx context.Context, senderPublicKey string, feeAssetId string) (uint64, error)

	CreateRawTxBySendersAddress(ctx context.Context, addressFrom string, addressTo string, amount uint64) ([]byte, error)
	CreateRawTxBySendersPublicKey(ctx context.Context, sendersPublicKey string, addressTo string, amount uint64, assetId string) ([]byte, error)
	SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error)
	GetTransactionByTxId(ctx context.Context, txId string) ([]byte, error)
	GetTransactionStatus(ctx context.Context, txId string) (models.TxStatus, error)
	TransactionByHash(ctx context.Context, txId string) (*models.TxInfo, error)
}

type nodeClient struct {
	conf config.Node
}

var (
	cl                    INodeClient
	onceRPCClientInstance sync.Once
)

// New create node's client with connection to Waves node
func New(ctx context.Context, conf config.Node) error {
	onceRPCClientInstance.Do(func() {
		cl = &nodeClient{conf: conf}
	})
	return nil
}

// GetNodeClient returns node's client.
// Client must be previously created with New(), in another case function throws panic
func GetNodeClient() INodeClient {
	onceRPCClientInstance.Do(func() {
		panic("try to get node client before it's creation!")
	})
	return cl
}
