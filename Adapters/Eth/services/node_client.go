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

	GetEthBalance(ctx context.Context, address string) (*big.Int, error)
	GetAllBalances(ctx context.Context, address string, contracts ...string) (*models.AccountBalance, error)
	GetNextNonce(ctx context.Context, address string) (uint64, error)
	GetErc20AllowanceAmount(ctx context.Context, ownerAddress string, contractAddress string, senderAddress string) (*big.Int, error)

	GenerateAddress(ctx context.Context) (publicAddress string, err error)
	IsAddressValid(ctx context.Context, address string) (bool, string, error)
	AddressByPublicKey(ctx context.Context, public string) (string, error)

	// create transaction and return it on RPL encoding
	CreateRawTransaction(ctx context.Context, addressFrom string, addressTo string, amount *big.Int, nonce uint64) ([]byte, error)
	CreateErc20TokensRawTransaction(ctx context.Context, addressFrom string, contractAddress string, addressTo string,
		amount *big.Int, nonce uint64) ([]byte, error)
	Erc20TokensRawApproveTransaction(ctx context.Context, ownerAddress string, contractAddress string, amount *big.Int,
		spenderAddress string) ([]byte, *big.Int, error)
	CreateErc20TokensTransferToTxSender(ctx context.Context, addressFrom string, contractAddress string,
		txSender string, amount *big.Int, nonce uint64) ([]byte, error)

	// sign transaction if has private key. rlpTx is transaction for signing on RPL encoding (returned by CreateRawTransaction function)
	SignTransaction(ctx context.Context, senderAddr string, rlpTx []byte) ([]byte, error)
	SignTransactionWithPrivateKey(ctx context.Context, privateKey string, rlpTx []byte) ([]byte, error)
	// send transaction. rlpTx is transaction for sending on RPL encoding (returned by SignTransaction function)
	SendTransaction(ctx context.Context, rlpTx []byte) (txHash string, err error)
	GetTxStatusByTxID(ctx context.Context, txID string) (models.TxStatus, error)
	TransactionInfo(ctx context.Context, txID string) (*models.TxInfo, error)
}

const gasLimitForMoneyTransfer = 21000

type nodeClient struct {
	ethClient        *ethclient.Client
	contractProvider *Erc20ContractProvider
	chainID          int64
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
		cp := NewContractProvider(ethClient)
		cl = &nodeClient{
			ethClient:        ethClient,
			contractProvider: cp,
			chainID:          config.ChainId,
			privateKeys:      make(map[string]*ecdsa.PrivateKey),
		}
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
