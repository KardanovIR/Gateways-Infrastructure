package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Ergo/config"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const httpRequestTimeoutMs = 2000

type INodeClient interface {
	BlockAt(ctx context.Context, blockNumber uint64, currentHeight *uint64) (*Block, error)
	BlockLast(ctx context.Context) (*BlockShortInfo, error)
}

type nodeClient struct {
	conf       config.Node
	httpClient http.Client
}

var (
	cl                    INodeClient
	onceRPCClientInstance sync.Once
)

// New create node's client with connection to ergo node
func NewNodeClient(ctx context.Context, conf config.Node) INodeClient {
	onceRPCClientInstance.Do(func() {
		tr := &http.Transport{}
		client := http.Client{
			Timeout:   time.Duration(httpRequestTimeoutMs) * time.Millisecond,
			Transport: tr,
		}
		cl = &nodeClient{conf: conf, httpClient: client}
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
	r, _ := cl.Request(ctx, http.MethodGet, cl.conf.Host+getCurrentBlockUrl, nil)
	getCurrentBlockResp := BlocksResponse{}
	if err := json.Unmarshal(r, &getCurrentBlockResp); err != nil {
		log.Errorf("failed to get current height: %s", err)
		return nil, err
	}
	block := getCurrentBlockResp.Items[0]
	log.Debugf("current height is %d", block.Height)
	return &block, nil
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
	blockInfo, err := cl.BlockByNumber(ctx, blockNumber, *currentHeight)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("blockInfo %+v ", blockInfo)
	log.Debugf("request for block with id %s", blockInfo.ID)
	r, _ := cl.Request(ctx, http.MethodGet, cl.conf.Host+fmt.Sprintf(getBlockByIdUrlTemplate, blockInfo.ID), nil)
	block := Block{}
	if err := json.Unmarshal(r, &block); err != nil {
		log.Errorf("failed to get block by id %s: %s", blockInfo.ID, err)
		return nil, err
	}
	log.Debugf("block %+v ", block)
	return &block, nil
}

func (cl *nodeClient) BlockByNumber(ctx context.Context, targetBlockNumber uint64, lastBlockHeight uint64) (*BlockShortInfo, error) {
	log := logger.FromContext(ctx)
	// calculate how much blocks will be skipped
	offset := lastBlockHeight - targetBlockNumber
	// request 2 blocks because height can be increased between requests
	r, _ := cl.Request(ctx, http.MethodGet, cl.conf.Host+fmt.Sprintf(getBlockByNumberUrlTemplate, offset), nil)
	blockResp := BlocksResponse{}
	if err := json.Unmarshal(r, &blockResp); err != nil {
		log.Errorf("failed to get current height: %s", err)
		return nil, err
	}
	// get 2 blocks, find block with needed block's height
	heights := make([]uint64, len(blockResp.Items))
	for i, bl := range blockResp.Items {
		heights[i] = bl.Height
		if bl.Height == targetBlockNumber {
			log.Debugf("BlockByNumber returns %d block", targetBlockNumber)
			return &bl, nil
		}
	}
	return nil, fmt.Errorf("haven't block with height %d between blocks %v", targetBlockNumber, heights)
}

func (cl *nodeClient) Request(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("wrong response status %d, body %s", resp.StatusCode, string(body))

	}
	respBody, err := ioutil.ReadAll(resp.Body)
	log := logger.FromContext(ctx)
	log.Debugf("response: %s", string(respBody))
	return respBody, err
}
