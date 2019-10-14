package data_client

import (
	"context"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
	"net/http"
	"sync"
	"time"
)

const httpRequestTimeoutMs = 2000

type IDataClient interface {
	GetAllBalances(ctx context.Context, address string) (*models.Balance, error)
	TransactionByHash(ctx context.Context, txId string) (*models.TxInfo, error)
	GetUnspentInputs(ctx context.Context, address string) ([]*models.RawUtxo, error)
}

type dataClient struct {
	conf       config.HttpService
	httpClient http.Client
}

var (
	dcl                    IDataClient
	onceHttpClientInstance sync.Once
)

// New create node's client with connection to Waves node
func NewDataClient(ctx context.Context, conf config.HttpService) error {
	onceHttpClientInstance.Do(func() {
		// configuration of TLS will be set here
		tr := &http.Transport{}
		client := http.Client{
			Timeout:   time.Duration(httpRequestTimeoutMs) * time.Millisecond,
			Transport: tr,
		}
		dcl = &dataClient{conf: conf, httpClient: client}
	})
	return nil
}

// GetNodeClient returns node's client.
// Client must be previously created with New(), in another case function throws panic
func GetDataClient() IDataClient {
	onceHttpClientInstance.Do(func() {
		panic("try to get node client before it's creation!")
	})
	return dcl
}
