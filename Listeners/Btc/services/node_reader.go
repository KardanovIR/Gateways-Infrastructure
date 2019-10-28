package services

import (
	"context"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/repository"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/converter"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/repositories"
	coreServices "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/services"
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
	nodeClient       INodeClient
	unspentTxService *unspentTxService
	rp               repositories.IRepository
	callbackService  coreServices.ICallbackService
	conf             *config.Node
	stopListen       chan struct{}
}

var (
	reader         INodeReader
	onceNodeClient sync.Once
)

func NewReader(ctx context.Context, config *config.Node, rp repository.IUTXORepository, nodeClient INodeClient,
	cs coreServices.ICallbackService) error {
	onceNodeClient.Do(func() {

		reader = &nodeReader{nodeClient: nodeClient,
			unspentTxService: NewUnspentTxService(ctx, rp),
			rp:               rp,
			conf:             config,
			callbackService:  cs,
			stopListen:       make(chan struct{}),
		}
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

	startBlockHeight := nr.conf.StartBlockHeight
	currentHeight, err := nr.nodeClient.GetCurrentHeight(ctx)
	if err != nil {
		log.Errorf("Current height error: %s", err)
		return err
	}
	if uint64(startBlockHeight) > currentHeight {
		log.Errorf("Configuration start block error: start block %d, current block %d.", startBlockHeight, currentHeight)
		startBlockHeight = currentHeight
	}
	log.Infof("Start listening Btc from %d block.", startBlockHeight)
	go func() {
		for {
			select {
			case <-nr.stopListen:
				log.Infof("Stop listening Btc.")
				return
			default:
			}

			log.Debugf("Process Btc block %d", startBlockHeight)
			log.Infof("Current block (%d), start block %d", currentHeight, startBlockHeight)
			if startBlockHeight > (currentHeight - nr.conf.Confirmations) {
				log.Debugf("confirmations of %d < %d. Waiting a %d minutes...", startBlockHeight, nr.conf.Confirmations, delayForWaitingNewBlock/time.Minute)
				time.Sleep(delayForWaitingNewBlock)
				nodeHeight, err := nr.nodeClient.GetCurrentHeight(ctx)
				if err != nil {
					log.Errorf("height error: %s", err)
					time.Sleep(delayAfterError)
				} else {
					currentHeight = nodeHeight
				}
				continue
			}

			block, err := nr.nodeClient.BlockAt(ctx, startBlockHeight)
			if err != nil {
				log.Errorf("BlockByNumber(%d) error: %s", startBlockHeight, err)
				time.Sleep(delayAfterError)
				continue
			}

			err = nr.processBlock(ctx, block)
			if err != nil {
				log.Errorf("processBlock(%d) error: %s. Waiting %d second", startBlockHeight, err, delayAfterError/time.Second)
				time.Sleep(delayAfterError)
				continue
			}

			if chainState == nil {
				chainState = new(models.ChainState)
				chainState.ChainType = models.ChainType(nr.conf.ChainType)
			}

			chainState.LastBlock = int64(startBlockHeight + 1)
			chainState.Timestamp = time.Now()

			chainState, err = nr.rp.PutChainState(ctx, chainState)
			if err != nil {
				log.Errorf("Updating chainState for %s error: %s", nr.conf.ChainType, err)
			}

			startBlockHeight = uint64(chainState.LastBlock)
		}
	}()

	return
}

func (nr *nodeReader) Stop(ctx context.Context) {
	log := logger.FromContext(ctx)
	log.Info("Stop listening Btc.")
	nr.stopListen <- struct{}{}
	return
}

func (nr *nodeReader) processBlock(ctx context.Context, block *GetBlockVerboseResult) (err error) {
	log := logger.FromContext(ctx)
	log.Debugf("start process block %d", block.Height)
	for _, tx := range block.Tx {
		if err := nr.processTx(ctx, tx); err != nil {
			return err
		}
	}
	return
}

func (nr *nodeReader) processTx(ctx context.Context, tx btcjson.TxRawResult) error {
	log := logger.FromContext(ctx)
	// find tasks for txId
	tasksForTxId, err := nr.rp.FindByAddressOrTxId(ctx, models.ChainType(nr.conf.ChainType), "", tx.Txid)
	if err != nil {
		log.Errorf("error: %s", err)
		return err
	}
	// if tasks for txId were found - it's our transaction for sending bitcoins - tx inputs were used - delete used inputs
	if len(tasksForTxId) > 0 {
		// get inputs for tx -> get inputTxId
		// delete inputs
		if err := nr.deleteTxInputs(ctx, tx.Vin); err != nil {
			return err
		}
	}

	tasks := make([]*models.Task, 0)
	for _, output := range tx.Vout {
		if len(output.ScriptPubKey.Addresses) != 1 {
			// many addresses can be for multisig account
			continue
		}
		address := output.ScriptPubKey.Addresses[0]
		tasksForAddress, err := nr.rp.FindByAddressOrTxId(ctx, models.ChainType(nr.conf.ChainType), address, "")
		if err != nil {
			log.Errorf("error: %s", err)
			return err
		}
		// if tasks for address were found - we interested in counting inputs for this address
		if len(tasksForAddress) > 0 {
			amount, err := converter.GetIntFromFloat(ctx, output.Value)
			if err != nil {
				log.Errorf("get amount fails: %s", err)
				return err
			}
			// add output to unspent inputs for address
			if err := nr.unspentTxService.addTxInputs(ctx, tx.Txid, amount, address, output.N); err != nil {
				return err
			}
			tasks = append(tasks, tasksForAddress...)
		}
	}
	tasks = append(tasks, tasksForTxId...)
	if len(tasks) > 0 {
		// execute tasks which have callbackUrl
		if err := nr.executeTasks(ctx, tasks, tx.Txid); err != nil {
			return err
		}
	}
	return nil
}

func (nr *nodeReader) deleteTxInputs(ctx context.Context, vIn []btcjson.Vin) error {
	// delete input
	// find UnspentTx by txId (of input) -> delete txId, update UnspentTx
	for _, in := range vIn {
		if err := nr.unspentTxService.deleteTxInputs(ctx, in.Txid, in.Vout); err != nil {
			return err
		}
	}
	return nil
}

func (nr *nodeReader) executeTasks(ctx context.Context, tasks []*models.Task, txId string) error {
	log := logger.FromContext(ctx)
	log.Infof("tx ID %s has %d tasks, start processing...", txId, len(tasks))
	var err error
	for _, task := range tasks {
		log.Infof("Start processing task id %s for %v ...", task.Id.Hex(), task.ListenTo)
		// skip tasks which was added only for calculation inputs for address and is not need callback call
		if len(task.Callback.Type) == 0 {
			log.Debugf("callback is empty for task %+v", task.ListenTo)
			continue
		}
		err = nr.callbackService.SendRequest(ctx, task, txId)
		if err != nil {
			log.Errorf("Error while processing task for %s %s: %s", task.ListenTo.Type, task.ListenTo.Value, err)
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
