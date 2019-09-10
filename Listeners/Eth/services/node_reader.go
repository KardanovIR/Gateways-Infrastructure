package services

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/repositories"
	coreServices "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/services"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/config"
)

type INodeReader interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context)
}

type nodeReader struct {
	nodeClient *nodeClient
	rp         repositories.IRepository
	conf       *config.Node
	stopListen chan struct{}
}

const FailedTxStatus = 0

var (
	cl             INodeReader
	onceNodeClient sync.Once
)

// New create node's client with connection to eth node
func New(ctx context.Context, config *config.Node, rp repositories.IRepository) error {
	log := logger.FromContext(ctx)
	var err error
	onceNodeClient.Do(func() {
		nc, e := newNodeClient(ctx, config.Host)
		if e != nil {
			err = e
		}
		cl = &nodeReader{nodeClient: nc, rp: rp, conf: config, stopListen: make(chan struct{})}
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

func (nr *nodeReader) Start(ctx context.Context) (err error) {
	log := logger.FromContext(ctx)
	client := nr.nodeClient

	confirmations := new(big.Int)
	confirmations.SetString(nr.conf.Confirmations, 10)

	//to do add saving chainState
	//example below
	chainState, errorChainState := nr.rp.GetLastChainState(ctx, models.ChainType(nr.conf.ChainType))
	if errorChainState == nil && chainState != nil {
		if (*chainState).LastBlock > nr.conf.StartBlockHeight {
			nr.conf.StartBlockHeight = chainState.LastBlock
		}
	}

	startBlock := big.NewInt(nr.conf.StartBlockHeight)

	log.Infof("Start listening ETH from %d block.", startBlock)
	go func() {
		for {
			select {
			case <-nr.stopListen:
				log.Infof("Stop listening ETH.")
				return
			default:
			}

			log.Infof("Process ETH block %d", startBlock)

			header, err := client.HeaderByNumber(ctx, startBlock)
			if err != nil {
				log.Errorf("BlockByNumber(%d) error: %s", startBlock, err)
				time.Sleep(15 * time.Second)
				continue
			}

			currentHeader, err := client.HeaderByNumber(ctx, nil)
			if err != nil {
				log.Errorf("HeaderByNumber error: %s", err)
				time.Sleep(15 * time.Second)
				continue
			}

			difference := currentHeader.Number
			log.Infof("Current block (%d), start block %d", difference, startBlock)
			difference.Sub(difference, startBlock)
			var compare = difference.Cmp(confirmations)
			if compare < 0 {
				log.Infof("Confirmations of %d < %d. Waiting a minute...", startBlock, nr.conf.Confirmations)
				time.Sleep(time.Minute)
				continue
			}

			block, err := client.BlockByHash(ctx, header.Hash())
			if err != nil {
				log.Errorf("BlockByHash error: %s", err)
				time.Sleep(15 * time.Second)
				continue
			}

			err = nr.processBlock(ctx, block)
			if err != nil {
				log.Errorf("processBTCBlock(%s) error: %s", startBlock, err)
				log.Debugf("Waiting a half minute.")
				time.Sleep(30 * time.Second)
				continue
			}

			if chainState == nil {
				chainState = new(models.ChainState)
				chainState.ChainType = models.ChainType(nr.conf.ChainType)
			}

			chainState.LastBlock = block.Number().Int64() + 1
			chainState.Timestamp = time.Now()

			chainState, err = nr.rp.PutChainState(ctx, chainState)
			if err != nil {
				log.Errorf("Updating chainState for %s error: %s", nr.conf.ChainType, err)
			}

			startBlock = big.NewInt(chainState.LastBlock)
		}

	}()

	return
}

func (nr *nodeReader) Stop(ctx context.Context) {
	log := logger.FromContext(ctx)
	log.Info("Stop listening ETH.")

	nr.stopListen <- struct{}{}
	return
}

func (nr *nodeReader) processBlock(ctx context.Context, block *types.Block) (err error) {
	log := logger.FromContext(ctx)

	log.Debugf("Start processing transactions of block %d...", block.Number())
	txs := block.Transactions()

	for _, tx := range txs {
		outAddresses := tx.To()
		txHash := tx.Hash().String()
		if outAddresses == nil {
			// contract creation transaction
			log.Debugf("nil address for tx %s", txHash)
			continue
		}

		isERC20Transfers, err := CheckERC20Transfers(tx.Data())
		if err != nil {
			log.Errorf("error checking erc20 tx", err)
		}
		if isERC20Transfers {
			receipt, err := nr.nodeClient.TransactionReceipt(ctx, tx.Hash())
			if err != nil {
				log.Errorf("get transaction receipt failed", err)
				return err
			}
			if receipt.Status == FailedTxStatus {
				log.Debugf("failed tx %s. skip it", txHash)
				continue
			}
			transferParams, err := ParseERC20TransferParams(tx.Data())
			if err != nil {
				log.Errorf("can't parse tx:", err)
				continue
			}
			outAddresses = &transferParams.To
		} else if len(tx.Data()) > 0 {
			// if data is not empty -> maybe it is call of the contract -> parse this with trace to find internal tx
			recipients, err := nr.nodeClient.GetEthRecipientsForTxIncludeInternal(ctx, txHash)
			if err != nil {
				log.Errorf("get eth internal transfer recipients fails: %s", err)
			} else {
				for _, r := range recipients {
					if err := nr.findAndExecuteTasks(ctx, r, txHash); err != nil {
						return err
					}
				}
			}
		}
		if err := nr.findAndExecuteTasks(ctx, outAddresses.Hex(), txHash); err != nil {
			return err
		}
	}
	return
}

func (nr *nodeReader) findAndExecuteTasks(ctx context.Context, address string, txHash string) error {
	log := logger.FromContext(ctx)
	tasks, err := nr.rp.FindByAddressOrTxId(ctx, models.ChainType(nr.conf.ChainType), address, txHash)
	if err != nil {
		log.Errorf("error: %s", err)
		return err
	}

	if len(tasks) == 0 {
		log.Debugf("tx %s has no task, will be skipped.", txHash)
		return nil
	}

	log.Debugf("tx %s has %d tasks, start processing...", txHash, len(tasks))

	for _, task := range tasks {
		log.Infof("Start processing task id %s for %v ...", task.Id.Hex(), task.ListenTo)
		err = coreServices.GetCallbackService().SendRequest(ctx, task, txHash)
		if err != nil {
			log.Errorf("Error: send callback %s for task %v for tx %s: %s", task.Callback.Type,
				task.ListenTo, txHash, err)
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
	return nil
}
