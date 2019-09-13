package services

import (
	"context"
	"github.com/BANKEX/payment-gateway-btc-adapter/_vendor-20180717124023/github.com/btcsuite/btcd/rpcclient"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
	"strconv"
	"sync"
)


type INodeClient interface {
	ValidateAddress(ctx context.Context, address string) (bool, error)
	PublicKeyFromAddress(ctx context.Context, address string) []byte
	GetAllBalances(ctx context.Context, address string) (*models.AccountBalance, error)

	Fee(ctx context.Context, senderPublicKey string) (uint64, error)

	CreateRawTx(ctx context.Context, addressFrom string, outs []*models.Output) ([]byte, error)
	SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error)
	TransactionByHash(ctx context.Context, txId string) (*models.TxInfo, error)
}

type nodeClient struct {
	conf       config.Node
	nodeClient *rpcclient.Client
}

var (
	cl                    INodeClient
	onceRPCClientInstance sync.Once
)

// New create node's client with connection to Waves node
func New(ctx context.Context, conf config.Node) error {
	onceRPCClientInstance.Do(func() {
		log := logger.FromContext(ctx)
		nodeCon := &rpcclient.ConnConfig{
			Host: conf.Host,
			User: conf.User,
			Pass: conf.Password,
			HTTPPostMode: conf.HTTPPostMode,
			DisableTLS: conf.DisableTLS,
		}

		client, err  := rpcclient.New(nodeCon, nil)
		if err != nil {
			log.Error(err)
		}
		cl = &nodeClient{conf: conf, nodeClient: client}
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