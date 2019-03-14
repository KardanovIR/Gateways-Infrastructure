package services

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/models"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
)

type INodeClient interface {
	GenerateAddress(ctx context.Context) (publicAddress string, err error)
	ValidateAddress(ctx context.Context, address string) (bool, error)
	GetLastBlockHeight(ctx context.Context) (string, error)
}

type nodeClient struct {
	nodeClient *client.Client
	chainID    models.NetworkType
	// private keys for addresses
	privateKeys map[string]crypto.SecretKey
}

func (cl *nodeClient) GetLastBlockHeight(ctx context.Context) (string, error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'GetLastBlockHeight'")

	lastBlock, _, err := cl.nodeClient.Blocks.Last(ctx)
	if err != nil {
		log.Errorf("get last block fails: %s", err)
		return "", err
	}
	return strconv.FormatUint(lastBlock.Height, 10), nil
}

var (
	cl                    INodeClient
	onceRPCClientInstance sync.Once
)

// New create node's client with connection to Waves node
func New(ctx context.Context, conf config.Node) error {
	log := logger.FromContext(ctx)
	var err error
	onceRPCClientInstance.Do(func() {
		wavesClient, e := client.NewClient(client.Options{
			Client:  &http.Client{Timeout: 30 * time.Second},
			BaseUrl: conf.Host,
		})
		if e != nil {
			log.Errorf("error during initialise waves client: %s", e)
			err = e
			return
		}

		cl = &nodeClient{nodeClient: wavesClient, chainID: conf.ChainID, privateKeys: make(map[string]crypto.SecretKey)}
	})
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
