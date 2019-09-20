package services

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"sync"

	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/config"
)

type INodeClient interface {
	BlockAt(ctx context.Context, blockNumber uint64, currentHeight *uint64) (*Block, error)
	BlockLast(ctx context.Context) (*BlockShortInfo, error)
}

type nodeClient struct {
	conf       config.Node
	nodeClient *rpcclient.Client
}

var (
	cl                    INodeClient
	onceRPCClientInstance sync.Once
)

// New create node's client with connection to btc node
func NewNodeClient(ctx context.Context, conf config.Node) INodeClient {
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
	return cl
}

// GetNodeClient returns node's client.
// Client must be previously created with NewNodeClient(), in another case function throws panic
func GetNodeClient() INodeClient {
	onceRPCClientInstance.Do(func() {
		panic("try to get node client before it's creation!")
	})
	return cl
}

const (
	getCurrentBlockUrl          = "/blocks?limit=1"
	getBlockByNumberUrlTemplate = "/blocks?limit=2&offset=%d"
	getBlockByIdUrlTemplate     = "/blocks/%s"
)

type BlocksResponse struct {
	Items []BlockShortInfo `json:"items"`
}

type BlockShortInfo struct {
	ID                string `json:"id"`
	Height            uint64 `json:"height"`
	TransactionsCount uint   `json:"transactionsCount"`
}

type Block struct {
	Block BlockEntity `json:"block"`
}

func (b *Block) Height() uint64 {
	return b.Block.Header.Height
}

func (b *Block) Transactions() []Transaction {
	return b.Block.Transactions
}

type BlockEntity struct {
	Header       Header        `json:"header"`
	Transactions []Transaction `json:"blockTransactions"`
}

type Header struct {
	Height uint64 `json:"height"`
}

type Transaction struct {
	ID        string      `json:"id"`
	TxOutputs []*TxOutput `json:"outputs"`
}

type TxOutput struct {
	Address string `json:"address"`
	Value   uint64 `json:"value"`
}

func (cl *nodeClient) BlockLast(ctx context.Context) (*BlockShortInfo, error) {
	log := logger.FromContext(ctx)
	log.Debug("get current height")
	//todo
	return nil, nil
}

func (cl *nodeClient) BlockAt(ctx context.Context, blockNumber uint64, currentHeight *uint64) (*Block, error) {
	log := logger.FromContext(ctx)
	log.Debugf("get BlockAt %d", blockNumber)
	// get current block
	if currentHeight == nil {
		lastBlock, err := cl.BlockLast(ctx)
		if err != nil {
			log.Errorf("failed to get current block: %s", err)
			return nil, err
		}
		currentHeight = &lastBlock.Height
	}
	//todo
	return nil, nil
}

func (cl *nodeClient) BlockByNumber(ctx context.Context, targetBlockNumber uint64, lastBlockHeight uint64) (*BlockShortInfo, error) {
	//log := logger.FromContext(ctx)
	// calculate how much blocks will be skipped
	//todo
	return nil, fmt.Errorf("haven't block with height %d between blocks %v", targetBlockNumber, 0)
}