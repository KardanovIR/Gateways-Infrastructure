package services

import (
	"context"
	"sync"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/repositories"
)

type INodeClient interface {
	ValidateAddress(ctx context.Context, address string) (bool, error)
	GetAllBalances(ctx context.Context, address string) (*models.Balance, error)
	GetBalanceForAllAddresses(ctx context.Context) (uint64, error)

	Fee(ctx context.Context) (uint64, error)

	CreateRawTx(ctx context.Context, addressesFrom []string, changeAddress string, outs []*models.Output) ([]byte, error)
	SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error)
	TransactionByHash(ctx context.Context, txId string) (*models.TxInfo, error)
	SignTransaction(ctx context.Context, txUnsigned []byte, privateKeyForAddress map[string]string) (txSigned []byte, err error)
}

type nodeClient struct {
	conf       config.Node
	nodeClient *rpcclient.Client
	rep        repositories.IRepository
}

var (
	cl                    INodeClient
	onceRPCClientInstance sync.Once
)

// New create node's client with connection to Btc node
func New(ctx context.Context, conf config.Node, rep repositories.IRepository) error {
	onceRPCClientInstance.Do(func() {
		log := logger.FromContext(ctx)
		nodeCon := &rpcclient.ConnConfig{
			Host:         conf.Host,
			User:         conf.User,
			Pass:         conf.Password,
			HTTPPostMode: conf.HTTPPostMode,
			DisableTLS:   conf.DisableTLS,
		}
		client, err := rpcclient.New(nodeCon, nil)
		if err != nil {
			log.Error(err)
		}
		cl = &nodeClient{conf: conf, nodeClient: client, rep: rep}
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
