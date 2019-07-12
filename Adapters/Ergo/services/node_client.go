package services

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/models"
)

const httpRequestTimeoutMs = 2000

type INodeClient interface {
	ValidateAddress(ctx context.Context, address string) (bool, error)
	PublicKeyFromAddress(ctx context.Context, address string) []byte
	GetAllBalances(ctx context.Context, address string) (*models.AccountBalance, error)

	Fee(ctx context.Context, senderPublicKey string, feeAssetId string) (uint64, error)

	CreateRawTx(ctx context.Context, addressFrom string, outs []*models.Output) ([]byte, error)
	SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error)
	TransactionByHash(ctx context.Context, txId string) (*models.TxInfo, error)
}

type nodeClient struct {
	conf       config.Node
	httpClient http.Client
}

var (
	cl                    INodeClient
	onceRPCClientInstance sync.Once
)

// New create node's client with connection to Waves node
func New(ctx context.Context, conf config.Node) error {
	onceRPCClientInstance.Do(func() {
		// configuration of TLS will be set here
		tr := &http.Transport{}
		client := http.Client{
			Timeout:   time.Duration(httpRequestTimeoutMs) * time.Millisecond,
			Transport: tr,
		}
		cl = &nodeClient{conf: conf, httpClient: client}
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
