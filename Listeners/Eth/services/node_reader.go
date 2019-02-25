package services

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/repositories"
	"sync"
)

type INodeReader interface {
	Start() (err error)
	Stop()
}

type nodeReader struct {
	nodeClient    *ethclient.Client
	rp            repositories.IRepository
	conf          *config.Node
	stopListenBTC chan struct{}
}

var (
	cl             INodeReader
	onceNodeClient sync.Once
)

// New create node's client with connection to eth node
func New(ctx context.Context, config config.Node, r repositories.IRepository) error {
	log := logger.FromContext(ctx)
	var err error
	onceNodeClient.Do(func() {
		rc, e := newRPCClient(log, config.Host)
		if e != nil {
			err = e
			return
		}
		ethClient := ethclient.NewClient(rc)
		cl = &nodeReader{nodeClient: ethClient, rp: r, conf: &config}
	})
	log.Errorf("error during initialise node client: %s", err)
	return err
}

// GetNodeReader returns node's reader instance.
// Client must be previously created with New(), in another case function throws panic
func GetNodeReader() INodeReader {
	onceNodeClient.Do(func() {
		panic("try to get node reader before it's creation!")
	})
	return cl
}

func newRPCClient(log logger.ILogger, host string) (*rpc.Client, error) {
	log.Infof("try to connect to etherium node", host)
	c, err := rpc.DialContext(context.Background(), host)
	if err != nil {
		log.Errorf("connect to etherium node fails: ", err)
		return nil, err
	}
	log.Info("connected to etherium node successfully")
	return c, nil
}

func (service *nodeReader) Start() (err error) {
	//client := service.nodeClient

	//startBlock := service.conf.StartBlockHeight

	return
}

func (service *nodeReader) Stop() {
	service.stopListenBTC <- struct{}{}
	return
}
