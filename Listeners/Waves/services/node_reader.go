package services

import (
	"context"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/repositories"
)

type INodeReader interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context)
}

type nodeReader struct {
	nodeClient *client.Client
	rp         repositories.IRepository
	rc         IRestClient
	conf       *config.Node
	stopListen chan struct{}
}

var (
	cl             INodeReader
	onceNodeClient sync.Once
)

// New create node's client with connection to eth node
func New(ctx context.Context, config *config.Node, rc IRestClient, rp repositories.IRepository) error {
	log := logger.FromContext(ctx)
	var err error
	onceNodeClient.Do(func() {
		wavesClient, e := client.NewClient(client.Options{
			Client:  &http.Client{Timeout: 30 * time.Second},
			BaseUrl: config.Host,
		})

		if e != nil {
			log.Errorf("error during initialise rpc client: %s", e)
			err = e
			return
		}

		cl = &nodeReader{nodeClient: wavesClient, rc: rc, rp: rp, conf: config}
	})

	if err != nil {
		log.Errorf("error during initialise node client: %s", err)
		return err
	}

	return nil
}

// GetNodeReader returns node's reader instance.
// Client must be previously created with New(), in another case function throws panic
func GetNodeReader() INodeReader {
	onceNodeClient.Do(func() {
		panic("try to get node reader before it's creation!")
	})
	return cl
}

func (service *nodeReader) Start(ctx context.Context) (err error) {
	log := logger.FromContext(ctx)
	client := service.nodeClient

	confirmations, err := strconv.ParseUint(service.conf.Confirmations, 10, 64)
	if err != nil {
		log.Errorf("Can't parse confirmations error: %s", err)
	}

	//to do add saving chainState
	//example below
	chainState, errorChainState := service.rp.GetLastChainState(ctx, models.ChainType(service.conf.ChainType))
	if errorChainState == nil && chainState != nil {
		if (*chainState).LastBlock > service.conf.StartBlockHeight {
			service.conf.StartBlockHeight = chainState.LastBlock
		}
	}

	startBlock := big.NewInt(service.conf.StartBlockHeight)

	log.Infof("Start listening ETH from %d block.", startBlock)
	go func() {
		for {
			select {
			case <-service.stopListen:
				log.Infof("Stop listening ETH.")
				return
			default:
			}

			log.Infof("Process ETH block %d", startBlock)

			block, _, err := client.Blocks.At(ctx, startBlock.Uint64())
			if err != nil {
				log.Errorf("BlockByNumber(%d) error: %s", startBlock, err)
				time.Sleep(15 * time.Second)
				continue
			}

			lastBlock, _, err := client.Blocks.Last(ctx)
			if err != nil {
				log.Errorf("HeaderByNumber error: %s", err)
				time.Sleep(15 * time.Second)
				continue
			}

			log.Infof("Current block (%d), start block %d", lastBlock.Height, block.Height)
			if block.Height > (lastBlock.Height - confirmations) {
				log.Infof("Confirmations of %d < %d. Waiting a minute...", startBlock, service.conf.Confirmations)
				time.Sleep(time.Minute)
				continue
			}

			err = service.processBlock(ctx, block)
			if err != nil {
				log.Errorf("processBTCBlock(%s) error: %s", startBlock, err)
				log.Debugf("Waiting a half minute.")
				time.Sleep(30 * time.Second)
				continue
			}

			if chainState == nil {
				chainState = new(models.ChainState)
				chainState.ChainType = models.ChainType(service.conf.ChainType)
			}

			chainState.LastBlock = int64(block.Height + 1)
			chainState.Timestamp = time.Now()

			*chainState, err = service.rp.PutChainState(ctx, *chainState)
			if err != nil {
				log.Errorf("Updating chainState for %s error: %s", service.conf.ChainType, err)
			}

			startBlock = big.NewInt(chainState.LastBlock)
		}

	}()

	return
}

func (service *nodeReader) Stop(ctx context.Context) {
	log := logger.FromContext(ctx)
	log.Info("Stop listening ETH.")

	service.stopListen <- struct{}{}
	return
}

func (service *nodeReader) processBlock(ctx context.Context, block *client.Block) (err error) {
	log := logger.FromContext(ctx)

	log.Infof("Start processing transactions of block %d...", block.Height)
	txs := block.Transactions


	for i, rawTx := range []proto.Transaction(txs) {
		//var scheme byte
		//var address crypto.PublicKey
		switch t := rawTx.(type) {
		case *proto.TransferV1:
			log.Debugf("%d: TransferV1: %v", i, t)
			tt := *t
			if tt.AmountAsset.Present || tt.FeeAsset.Present {

				//u, err := data.FromTransferV1(scheme, tt, address)
				//if err != nil {
				//	return err
				//}
				//_ = u
			}

		case *proto.TransferV2:
			log.Debugf("%d: TransferV2: %v", i, t)
			tt := *t
			if tt.AmountAsset.Present || tt.FeeAsset.Present {
				//u, err := data.FromTransferV2(scheme, tt, address)
				//if err != nil {
				//	return err
				//}
				//_=u
			}
		}



	}

	return
}
