package services

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/models"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type INodeClient interface {
	GenerateAddress(ctx context.Context) (publicAddress string, err error)
	ValidateAddress(ctx context.Context, address string) (bool, error)
	GetBalance(ctx context.Context, address string) (uint64, error)

	Fee(ctx context.Context, senderPublicKey string, feeAssetId string) (uint64, error)
	FeeForTx(ctx context.Context, tx *proto.TransferV2) (uint64, error)
	GetLastBlockHeight(ctx context.Context) (string, error)

	CreateRawTxBySendersAddress(ctx context.Context, addressFrom string, addressTo string, amount uint64) ([]byte, error)
	CreateRawTxBySendersPublicKey(ctx context.Context, sendersPublicKey string, addressTo string, amount uint64, assetId string) ([]byte, error)
	SignTxWithKeepedSecretKey(ctx context.Context, sendersAddress string, txUnsigned []byte) ([]byte, error)
	SignTxWithSecretKey(ctx context.Context, secretKeyInBase58 string, txUnsigned []byte) ([]byte, error)
	SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error)
	GetTransactionByTxId(ctx context.Context, txId string) ([]byte, error)
	GetTransactionStatus(ctx context.Context, txId string) (models.TxStatus, error)
}

type nodeClient struct {
	nodeClient *client.Client
	chainID    models.NetworkType
	// private keys for addresses
	privateKeys map[string]crypto.SecretKey
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
