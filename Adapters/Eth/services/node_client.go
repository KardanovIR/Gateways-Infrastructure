package services

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/models"
)

type INodeClient interface {
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	SuggestFee(ctx context.Context) (*big.Int, error)

	GetBalance(ctx context.Context, address string) (*big.Int, error)
	GetNextNonce(ctx context.Context, address string) (uint64, error)

	GenerateAddress(ctx context.Context) (publicAddress string, err error)
	IsAddressValid(ctx context.Context, address string) bool

	// create transaction and return it on RPL encoding
	CreateRawTransaction(ctx context.Context, addressFrom string, addressTo string, amount *big.Int) ([]byte, error)
	// sign transaction if has private key. rlpTx is transaction for signing on RPL encoding (returned by CreateRawTransaction function)
	SignTransaction(ctx context.Context, senderAddr string, rlpTx []byte) ([]byte, error)
	SignTransactionWithPrivateKey(ctx context.Context, privateKey string, rlpTx []byte) ([]byte, error)
	// send transaction. rlpTx is transaction for sending on RPL encoding (returned by SignTransaction function)
	SendTransaction(ctx context.Context, rlpTx []byte) (txHash string, err error)
	GetTxStatusByTxID(ctx context.Context, txID string) (models.TxStatus, error)
}

const gasLimitForMoneyTransfer = 21000

type nodeClient struct {
	ethClient *ethclient.Client
	chainID   int64
	// private keys for addresses
	privateKeys map[string]*ecdsa.PrivateKey
}

var (
	cl                    INodeClient
	onceRPCClientInstance sync.Once
)

// New create node's client with connection to eth node
func New(ctx context.Context, config config.Node) error {
	log := logger.FromContext(ctx)
	var err error
	onceRPCClientInstance.Do(func() {
		rc, e := newRPCClient(log, config.Host)
		if e != nil {
			err = e
			return
		}
		ethClient := ethclient.NewClient(rc)
		cl = &nodeClient{ethClient: ethClient, chainID: config.ChainId, privateKeys: make(map[string]*ecdsa.PrivateKey)}
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
