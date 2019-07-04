package services

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/repositories"
	coreServices "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/services"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Ergo/config"
)

const (
	delayAfterError         = time.Second * 30
	delayForWaitingNewBlock = time.Minute
)

type INodeReader interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context)
}

type nodeReader struct {
	nodeClient      INodeClient
	rp              repositories.IRepository
	callbackService coreServices.ICallbackService
	conf            *config.Node
	stopListen      chan struct{}
}

var (
	reader         INodeReader
	onceNodeClient sync.Once
)

func NewReader(ctx context.Context, config *config.Node, rp repositories.IRepository, nodeClient INodeClient,
	cs coreServices.ICallbackService) error {
	onceNodeClient.Do(func() {
		reader = &nodeReader{nodeClient: nodeClient, rp: rp, conf: config, callbackService: cs, stopListen: make(chan struct{})}
	})
	return nil
}

// GetNodeReader returns node's reader instance.
// Client must be previously created with New(), in another case function throws panic
func GetNodeReader() INodeReader {
	onceNodeClient.Do(func() {
		panic("try to get node reader before it's creation!")
	})
	return reader
}

func (nr *nodeReader) Start(ctx context.Context) (err error) {
	log := logger.FromContext(ctx)

	chainState, errorChainState := nr.rp.GetLastChainState(ctx, models.ChainType(nr.conf.ChainType))
	if errorChainState == nil && chainState != nil {
		if uint64(chainState.LastBlock) > nr.conf.StartBlockHeight {
			nr.conf.StartBlockHeight = uint64(chainState.LastBlock)
		}
	}

	startBlock := nr.conf.StartBlockHeight
	lastBlock, err := nr.nodeClient.BlockLast(ctx)
	if err != nil {
		log.Errorf("Last block error: %s", err)
	}
	if uint64(startBlock) > lastBlock.Height {
		log.Errorf("Configuration start block error: start block %d, current block %d.", startBlock, lastBlock.Height)
		startBlock = lastBlock.Height
	}
	log.Infof("Start listening Ergo from %d block.", startBlock)
	go func() {
		for {
			select {
			case <-nr.stopListen:
				log.Infof("Stop listening Ergo.")
				return
			default:
			}

			log.Debugf("Process Ergo block %d", startBlock)
			lastBlock, err := nr.nodeClient.BlockLast(ctx)
			if err != nil {
				log.Errorf("last block error: %s", err)
				time.Sleep(delayAfterError)
				continue
			}

			log.Infof("Current block (%d), start block %d", lastBlock.Height, startBlock)
			if startBlock > (lastBlock.Height - nr.conf.Confirmations) {
				log.Debugf("confirmations of %d < %d. Waiting a %d minutes...", startBlock, nr.conf.Confirmations, delayForWaitingNewBlock/time.Minute)
				time.Sleep(delayForWaitingNewBlock)
				continue
			}

			block, err := nr.nodeClient.BlockAt(ctx, startBlock, &lastBlock.Height)
			if err != nil {
				log.Errorf("BlockByNumber(%d) error: %s", startBlock, err)
				time.Sleep(delayAfterError)
				continue
			}

			err = nr.processBlock(ctx, block)
			if err != nil {
				log.Errorf("processBlock(%d) error: %s. Waiting %d second", startBlock, err, delayAfterError/time.Second)
				time.Sleep(delayAfterError)
				continue
			}

			if chainState == nil {
				chainState = new(models.ChainState)
				chainState.ChainType = models.ChainType(nr.conf.ChainType)
			}

			chainState.LastBlock = int64(block.Height() + 1)
			chainState.Timestamp = time.Now()

			chainState, err = nr.rp.PutChainState(ctx, chainState)
			if err != nil {
				log.Errorf("Updating chainState for %s error: %s", nr.conf.ChainType, err)
			}

			startBlock = uint64(chainState.LastBlock)
		}
	}()

	return
}

func (nr *nodeReader) Stop(ctx context.Context) {
	log := logger.FromContext(ctx)
	log.Info("Stop listening Ergo.")
	nr.stopListen <- struct{}{}
	return
}

func (nr *nodeReader) processBlock(ctx context.Context, block *Block) (err error) {
	log := logger.FromContext(ctx)
	log.Debugf("start process block %d", block.Height())
	for _, tx := range block.Transactions() {
		for _, output := range tx.TxOutputs {
			address := strings.ToLower(output.Address)
			if err := nr.findAndExecuteTasks(ctx, address, tx.ID); err != nil {
				return err
			}
		}
	}
	return
}

func (nr *nodeReader) findAndExecuteTasks(ctx context.Context, address string, txId string) error {
	log := logger.FromContext(ctx)
	tasks, err := nr.rp.FindByAddressOrTxId(ctx, models.ChainType(nr.conf.ChainType), address, txId)
	if err != nil {
		log.Errorf("error: %s", err)
		return err
	}
	if len(tasks) == 0 {
		return nil
	}
	log.Infof("tx ID %s has %d tasks, start processing...", txId, len(tasks))
	for _, task := range tasks {
		log.Debugf("Start processing task id %s for %v ...", task.Id.Hex(), task.ListenTo)
		err = nr.callbackService.SendRequest(ctx, task, txId)
		if err != nil {
			log.Errorf("Error while processing tx %s for address %s: %s", txId, address, err)
			continue
		}
		switch task.Type {
		case models.OneTime:
			err := nr.rp.RemoveTask(ctx, task.Id.Hex())
			if err != nil {
				log.Errorf("Error: removing task %s", err)
				return err
			}
		}
		log.Debugf("Task id %s for %s has been proceed successful!", task.Id.Hex(), nr.conf.ChainType)
	}
	return err
}
