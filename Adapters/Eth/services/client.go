package services

import (
	"context"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type INodeClient interface {
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}

type nodeClient struct {
	ethClient *ethclient.Client
}

func (cl *nodeClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return cl.ethClient.SuggestGasPrice(ctx)
}

var (
	cl                    INodeClient
	onceRPCClientInstance sync.Once
)

// New create node's client with connection to eth node
func New(host string) error {
	var err error
	onceRPCClientInstance.Do(func() {
		rc, e := newRPCClient(host)
		if e != nil {
			err = e
			return
		}
		ethClient := ethclient.NewClient(rc)
		cl = &nodeClient{ethClient: ethClient}
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

func newRPCClient(host string) (*rpc.Client, error) {
	log.Println("try to connect to etherium node", host)
	c, err := rpc.DialContext(context.Background(), host)
	if err != nil {
		log.Println("connect to etherium node fails: ", err)
		return nil, err
	}
	log.Println("connected to etherium node successfully")
	return c, nil
}
