package services

import (
	"context"
	"errors"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/repositories"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type INodeReader interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context)
}

type nodeReader struct {
	nodeClient *client.Client
	rp         repositories.IRepository
	conf       *config.Node
	stopListen chan struct{}
}

var (
	cl             INodeReader
	onceNodeClient sync.Once
)

// New create node's client with connection to Waves node
func New(ctx context.Context, config *config.Node, rp repositories.IRepository) error {
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

		cl = &nodeReader{nodeClient: wavesClient, rp: rp, conf: config, stopListen: make(chan struct{})}
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

	log.Infof("Start listening Waves from %d block.", startBlock)
	go func() {
		for {
			select {
			case <-service.stopListen:
				log.Infof("Stop listening Waves.")
				return
			default:
			}

			log.Infof("Process Waves block %d", startBlock)

			block, _, err := client.Blocks.At(ctx, startBlock.Uint64())
			if err != nil {
				log.Errorf("BlockByNumber(%d) error: %s", startBlock, err)
				time.Sleep(30 * time.Second)
				continue
			}

			lastBlock, _, err := client.Blocks.Last(ctx)
			if err != nil {
				log.Errorf("Last block error: %s", err)
				time.Sleep(15 * time.Second)
				continue
			}

			log.Infof("Current block (%d), start block %d", lastBlock.Height, block.Height)
			if block.Height > (lastBlock.Height - confirmations) {
				log.Debugf("Confirmations of %d < %d. Waiting a minute...", startBlock, service.conf.Confirmations)
				time.Sleep(time.Minute)
				continue
			}

			err = service.processBlock(ctx, block)
			if err != nil {
				log.Errorf("processBlock(%s) error: %s", startBlock, err)
				log.Debug("Waiting a half minute.")
				time.Sleep(30 * time.Second)
				continue
			}

			if chainState == nil {
				chainState = new(models.ChainState)
				chainState.ChainType = models.ChainType(service.conf.ChainType)
			}

			chainState.LastBlock = int64(block.Height + 1)
			chainState.Timestamp = time.Now()

			chainState, err = service.rp.PutChainState(ctx, chainState)
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
	log.Info("Stop listening Waves.")

	service.stopListen <- struct{}{}
	return
}

func (service *nodeReader) processBlock(ctx context.Context, block *client.Block) (err error) {
	log := logger.FromContext(ctx)

	log.Infof("Start processing transactions of block %d...", block.Height)
	txs := block.Transactions
	for _, rawTx := range txs {
		switch t := rawTx.(type) {
		case *proto.TransferV1:
			log.Debugf("parse transaction type %d ('TransferV1'). ID %s", proto.TransferTransaction, t.ID)
			tt := *t
			var assetId crypto.Digest
			if tt.AmountAsset.Present {
				assetId = tt.AmountAsset.ID
			}
			if err := service.executeTasks(ctx, tt.Amount, tt.Recipient, assetId, t.ID.String()); err != nil {
				return err
			}
		case *proto.TransferV2:
			log.Debugf("parse transaction type %d ('TransferV2'). ID %s", proto.TransferTransaction, t.ID)
			tt := *t
			var assetId crypto.Digest
			if tt.AmountAsset.Present {
				assetId = tt.AmountAsset.ID
			}
			if err := service.executeTasks(ctx, tt.Amount, tt.Recipient, assetId, t.ID.String()); err != nil {
				return err
			}
		case *proto.MassTransferV1:
			log.Debugf("parse transaction type %d ('MassTransferTransaction'). ID %s", proto.MassTransferTransaction, t.ID.String())
			tt := *t
			var assetId crypto.Digest
			if tt.Asset.Present {
				assetId = tt.Asset.ID
			}
			for _, transfer := range tt.Transfers {
				// TxId is one for all transfers. But tx task is OneTime task (by logic)
				// That's why execution of first transfer tracks address and txId, txId task will be removed
				if err := service.executeTasks(ctx, transfer.Amount, transfer.Recipient, assetId, tt.ID.String()); err != nil {
					return err
				}
			}
		case *proto.Payment:
			// payment is only for waves transfer
			log.Debugf("parse transaction type %d ('Payment'). ID %s", proto.PaymentTransaction, t.ID)
			tt := *t
			if err := service.executeTasksForRecipientOrTxId(ctx, tt.Amount, tt.Recipient.String(), crypto.Digest{}, t.ID.String()); err != nil {
				return err
			}
		default:
			log.Debugf("not interesting transaction type for id %+v", t)
		}
	}
	return
}

func (service *nodeReader) executeTasks(ctx context.Context, amount uint64, recipient proto.Recipient,
	assetId crypto.Digest, txId string) error {
	if recipient.Address != nil {
		address := recipient.Address.String()
		return service.executeTasksForRecipientOrTxId(ctx, amount, address, assetId, txId)
	}
	if recipient.Alias != nil {
		alias := recipient.Alias.Alias
		return service.executeTasksForRecipientOrTxId(ctx, amount, alias, assetId, txId)
	}
	return errors.New("haven't recipient address")
}

func (service *nodeReader) executeTasksForRecipientOrTxId(ctx context.Context, amount uint64, recipient string,
	assetId crypto.Digest, txId string) (err error) {

	log := logger.FromContext(ctx)
	tasks, err := service.rp.FindByAddressOrTxId(ctx, models.ChainType(service.conf.ChainType), recipient, txId)
	if err != nil {
		log.Errorf("error: %s", err)
		return err
	}
	if len(tasks) == 0 {
		return nil
	}
	log.Infof("address %s has %d tasks, start processing...", recipient, len(tasks))
	for _, task := range tasks {
		log.Debugf("Start processing task id %s for %v ...", task.Id.Hex(), task.ListenTo)
		err = GetCallbackService().SendRequest(ctx, task, txId)
		if err != nil {
			log.Errorf("Error while processing incoming transfer for address %s. %s", recipient, err)
			continue
		}
		switch task.Type {
		case models.OneTime:
			err := service.rp.RemoveTask(ctx, task.Id.Hex())
			if err != nil {
				log.Errorf("Error: removing task %s", err)
				return err
			}
		}
		log.Debugf("Task id %s for %s has been proceed successful!", task.Id.Hex(), service.conf.ChainType)
	}
	return err
}
