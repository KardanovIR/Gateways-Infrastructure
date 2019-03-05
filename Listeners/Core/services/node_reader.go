package services

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/repositories"
)

type INodeReader interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context)
}

type nodeReader struct {
	nodeClient *ethclient.Client
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
		rpcc, e := newRPCClient(log, config.Host)
		if e != nil {
			log.Errorf("error during initialise rpc client: %s", e)
			err = e
			return
		}
		ethClient := ethclient.NewClient(rpcc)
		cl = &nodeReader{nodeClient: ethClient, rc: rc, rp: rp, conf: config}
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

func (service *nodeReader) Start(ctx context.Context) (err error) {
	log := logger.FromContext(ctx)
	client := service.nodeClient

	confirmations := new(big.Int)
	confirmations.SetString(service.conf.Confirmations, 10)

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
				log.Infof("Confirmations of %d < %d. Waiting a minute...", startBlock, service.conf.Confirmations)
				time.Sleep(time.Minute)
				continue
			}

			block, err := client.BlockByHash(ctx, header.Hash())
			if err != nil {
				log.Errorf("BlockByNumber(%d) error: %s", startBlock, err)
				time.Sleep(15 * time.Second)
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

			chainState.LastBlock = block.Number().Int64() + 1
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

func (service *nodeReader) processBlock(ctx context.Context, block *types.Block) (err error) {
	log := logger.FromContext(ctx)

	log.Infof("Start processing transactions of block %d...", block.Number())
	txs := block.Transactions()

	for _, tx := range txs {
		outAddresses := tx.To()

		if outAddresses == nil {
			log.Infof("nil address", err)
			continue
		}

		tasks, err := service.rp.FindByAddress(ctx, models.ChainType(service.conf.ChainType), outAddresses.Hex())
		if err != nil {
			log.Errorf("error: %s", err)
			return err
		}

		if len(tasks) == 0 {
			log.Infof("-> tx %s has no task, will be skipped.", tx.Hash().String())
			continue
		}

		log.Infof("-> tx %s has %d transafers, start processing...", tx.Hash().String(), len(tasks))

		for _, task := range tasks {
			log.Infof("->   Start processing transfer on registered wallet %s ...", task.Address)
			err := service.rc.RequestCallback(ctx, task.Callback, task.Callback.Data)
			//block, btcTx, wallet.Address, outAddresses[wallet.Address])
			if err != nil {
				log.Errorf("->   Error: creating incoming tx %s for wallet %s: %s", tx.Hash().String(), task.Address, err)
				continue
			}

			switch task.Type {
			case models.OneTime:
				err := service.rp.RemoveTask(ctx, string(task.Id))
				if err != nil {
					log.Errorf("->   Error: removing task %s", err)
					return err
				}
			}

			log.Infof("->   Transfer %s to %s has been registered!", service.conf.ChainType, task.Address)
		}

	}

	return
}
