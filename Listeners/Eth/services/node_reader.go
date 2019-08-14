package services

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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
	nodeClient *ethclient.Client
	rp         repositories.IRepository
	conf       *config.Node
	stopListen chan struct{}
}

var (
	cl             INodeReader
	onceNodeClient sync.Once
)

// New create node's client with connection to eth node
func New(ctx context.Context, config *config.Node, rp repositories.IRepository) error {
	log := logger.FromContext(ctx)
	var err error
	onceNodeClient.Do(func() {
		rpcc, e := newRPCClient(log, config.Host)
		if e != nil {
			log.Errorf("error during initialise rpc client: %s", e)
			err = e
			return
		}
		ethClient := ethclient.NewClient(rpcc)
		cl = &nodeReader{nodeClient: ethClient, rp: rp, conf: config, stopListen: make(chan struct{})}
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

func newRPCClient(log logger.ILogger, host string) (*rpc.Client, error) {
	log.Infof("try to connect to etherium node", host)
	c, err := rpc.DialContext(context.Background(), host)
	if err != nil {
		log.Errorf("connect to etherium node fails: ", err)
		return nil, err
	}
	log.Info("connected to etherium node successfully")
	return c, nil
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
				log.Errorf("BlockByNumber(%d) error: %s", startBlock, err)
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

	log.Infof("Start processing transactions of block %d...", block.Number())
	txs := block.Transactions()

	for _, tx := range txs {
		outAddresses := tx.To()

		if outAddresses == nil {
			log.Debugf("nil address:", err)
			continue
		}

		isERC20Transfers, err := CheckERC20Transfers(tx.Data())
		if err != nil {
			log.Debugf("error checking erc20 tx", err)
		}
		if isERC20Transfers {
			transferParams, err := ParseERC20TransferParams(tx.Data())
			if err != nil {
				log.Debugf("can't parse tx:", err)
				continue
			}
			outAddresses = &transferParams.To
		}

		tasks, err := nr.rp.FindByAddressOrTxId(ctx, models.ChainType(nr.conf.ChainType), outAddresses.Hex(), tx.Hash().Hex())
		if err != nil {
			log.Errorf("error: %s", err)
			return err
		}

		if len(tasks) == 0 {
			log.Debugf("-> tx %s has no task, will be skipped.", tx.Hash().String())
			continue
		}

		log.Debugf("-> tx %s has %d tasks, start processing...", tx.Hash().String(), len(tasks))

		for _, task := range tasks {
			log.Infof("->   Start processing task id %s for %v ...", task.Id.Hex(), task.ListenTo)
			err = coreServices.GetCallbackService().SendRequest(ctx, task, tx.Hash().String())
			if err != nil {
				log.Errorf("->   Error: send callback %s for task %v for tx %s: %s", task.Callback.Type,
					task.ListenTo, tx.Hash().String(), err)
				continue
			}

			switch task.Type {
			case models.OneTime:
				err := nr.rp.RemoveTask(ctx, task.Id.Hex())
				if err != nil {
					log.Errorf("->   Error: removing task %s", err)
					return err
				}
			}

			log.Debugf("->   Task id %s for %s has been proceed successful!", task.Id.Hex(), nr.conf.ChainType)
		}

	}

	return
}
