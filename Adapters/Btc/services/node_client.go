package services

import (
	"context"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
	"sync"
)

type INodeClient interface {
	ValidateAddress(ctx context.Context, address string) (bool, error)

	Fee(ctx context.Context, senderPublicKey string) (uint64, error)

	CreateRawTx(ctx context.Context, addressFrom string, outs []*models.Output) ([]byte, error)
	SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error)
}

type nodeClient struct {
	conf       config.Node
	nodeClient *rpcclient.Client
}

var (
	cl                    INodeClient
	onceRPCClientInstance sync.Once
)

// New create node's client with connection to Btc node
func New(ctx context.Context, conf config.Node) error {
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
