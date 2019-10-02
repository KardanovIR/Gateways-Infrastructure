package services

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
)

type INodeClient interface {
	BlockAt(ctx context.Context, blockNumber uint64) (*GetBlockVerboseResult, error)
	GetCurrentHeight(ctx context.Context) (uint64, error)
	GetBlockVerboseTx(ctx context.Context, blockHash string) (*GetBlockVerboseResult, error)
}

const (
	getBlockMethod           = "getblock"
	verbosityBlockWithTxType = 2
)

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

func (cl *nodeClient) GetCurrentHeight(ctx context.Context) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Debug("get current height")
	chainInfo, err := cl.nodeClient.GetBlockChainInfo()
	if err != nil {
		log.Errorf("get blockchain info fails: %s", err)
		return 0, err
	}
	return uint64(chainInfo.Blocks), nil
}

func (cl *nodeClient) BlockAt(ctx context.Context, blockNumber uint64) (*GetBlockVerboseResult, error) {
	log := logger.FromContext(ctx)
	log.Debugf("get BlockAt %d", blockNumber)
	// get current block
	blockHash, err := cl.nodeClient.GetBlockHash(int64(blockNumber))
	if err != nil {
		log.Errorf("failed to get block %d hash: %s", blockNumber, err)
		return nil, err
	}

	block, err := cl.GetBlockVerboseTx(ctx, blockHash.String())
	if err != nil {
		log.Errorf("failed to get block %d: %s", blockNumber, err)
		return nil, err
	}
	return block, nil
}

type GetBlockVerboseResult struct {
	Hash          string                `json:"hash"`
	Confirmations int64                 `json:"confirmations"`
	StrippedSize  int32                 `json:"strippedsize"`
	Size          int32                 `json:"size"`
	Weight        int32                 `json:"weight"`
	Height        int64                 `json:"height"`
	Version       int32                 `json:"version"`
	VersionHex    string                `json:"versionHex"`
	MerkleRoot    string                `json:"merkleroot"`
	Tx            []btcjson.TxRawResult `json:"tx,omitempty"`
	Time          int64                 `json:"time"`
	Nonce         uint32                `json:"nonce"`
	Bits          string                `json:"bits"`
	Difficulty    float64               `json:"difficulty"`
	PreviousHash  string                `json:"previousblockhash"`
	NextHash      string                `json:"nextblockhash,omitempty"`
}

func (cl *nodeClient) GetBlockVerboseTx(ctx context.Context, blockHash string) (*GetBlockVerboseResult, error) {
	log := logger.FromContext(ctx)
	blHash, err := json.Marshal(blockHash)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	verbosity, _ := json.Marshal(verbosityBlockWithTxType)
	rawResult, err := cl.nodeClient.RawRequest(getBlockMethod, []json.RawMessage{blHash, verbosity})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	blockVerboseRes := GetBlockVerboseResult{}
	if err = json.Unmarshal(rawResult, &blockVerboseRes); err != nil {
		log.Error(err)
		return nil, err
	}
	return &blockVerboseRes, err
}
