package services

import (
	"context"
	"log"
	"sync"

	"github.com/GatewaysInfrastructure/Listeners/Eth/repositories"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type INodeReader interface {
}

type nodeReader struct {
	nodeClient *ethclient.Client
	rp         repositories.IRepository
}

var (
	cl             INodeReader
	onceNodeClient sync.Once
)

// New create node's client with connection to eth node
func New(host string, r repositories.IRepository) error {
	var err error
	onceNodeClient.Do(func() {
		rc, e := newRPCClient(host)
		if e != nil {
			err = e
			return
		}
		ethClient := ethclient.NewClient(rc)
		cl = &nodeReader{nodeClient: ethClient, rp: r}
	})
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
