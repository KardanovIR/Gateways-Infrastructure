package services

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
)

type INodeClient interface {
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}

type nodeClient struct {
	ethClient *ethclient.Client
}

func (cl *nodeClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'SuggestGasPrice'")
	return cl.ethClient.SuggestGasPrice(ctx)
}

var (
	cl                    INodeClient
	onceRPCClientInstance sync.Once
)

// New create node's client with connection to eth node
func New(ctx context.Context, host string) error {
	log := logger.FromContext(ctx)
	var err error
	onceRPCClientInstance.Do(func() {
		rc, e := newRPCClient(log, host)
		if e != nil {
			err = e
			return
		}
		ethClient := ethclient.NewClient(rc)
		cl = &nodeClient{ethClient: ethClient}
	})
	log.Errorf("error during initialise node client: %s", err)
	return err
}

// GetNodeClient returns node's client.
// Client must be previously created with New(), in another case function throws panic
func GetNodeClient() INodeClient {
	onceRPCClientInstance.Do(func() {
		panic("try to get node client before it's creation!")
	})
	return cl
}

func newRPCClient(log logger.ILogger, host string) (*rpc.Client, error) {
	log.Infof("try to connect to etherium node %s", host)
	c, err := rpc.DialContext(context.Background(), host)
	if err != nil {
		log.Errorf("connect to etherium node fails: %s", err)
		return nil, err
	}
	log.Info("connected to etherium node successfully")
	return c, nil
}
