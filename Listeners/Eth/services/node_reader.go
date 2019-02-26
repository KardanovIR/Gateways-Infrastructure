package services

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/repositories"
	"math/big"
	"sync"
	"time"
)

type INodeReader interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context)
}

type nodeReader struct {
	nodeClient    *ethclient.Client
	rp            repositories.IRepository
	conf          *config.Node
	stopListenBTC chan struct{}
}

var (
	cl             INodeReader
	onceNodeClient sync.Once
)

// New create node's client with connection to eth node
func New(ctx context.Context, config config.Node, r repositories.IRepository) error {
	log := logger.FromContext(ctx)
	var err error
	onceNodeClient.Do(func() {
		rc, e := newRPCClient(log, config.Host)
		if e != nil {
			err = e
			return
		}
		ethClient := ethclient.NewClient(rc)
		cl = &nodeReader{nodeClient: ethClient, rp: r, conf: &config}
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

	startBlock := new(big.Int)
	startBlock.SetString(service.conf.StartBlockHeight, 10)

	confirmations := new(big.Int)
	confirmations.SetString(service.conf.Confirmations, 10)

	//to do add saving chainState
	//example below
	//chainState, errorChainState := service.dbConn.GetChainState(service.conf.Ticker)
	//if errorChainState == nil && chainState != nil {
	//	if (*chainState).LastBlock > startBlock {
	//		service.conf.BTCConfig.StartBlockHeight = chainState.LastBlock
	//	}
	//}

	log.Infof("Start listening ETH from %d block.", startBlock)
	go func() {
		for {
			select {
			case <-service.stopListenBTC:
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
				log.Infof("processBTCBlock(%s) error: %s", startBlock, err)
				log.Infof("Waiting a half minute.")
				time.Sleep(30 * time.Second)
				continue
			}
		}

	}()

	return
}

func (service *nodeReader) Stop(ctx context.Context) {
	log := logger.FromContext(ctx)
	log.Info("Stop listening ETH.")

	service.stopListenBTC <- struct{}{}
	return
}

func (service *nodeReader) processBlock(ctx context.Context, block *types.Block) (err error) {
	log := logger.FromContext(ctx)

	log.Infof("Start processing transactions of block %d...", block.Number())
	txs := block.Transactions()

	for _, tx := range txs {
		outAddresses := tx.To().Hex()

		task, err := service.rp.FindByAddress(ctx, service.conf.Ticker, outAddresses)
		if err != nil {
			return err
		}

		if len(task) == 0 {
			log.Infof("-> tx %s has no transfers, will be skipped.", tx.Hash())
			continue
		}

		log.Infof("-> tx %s has %d transafers, start processing...", tx.Hash(), len(task))

	//	for _, wallet := range wallets {
	//		log.Printf("->   Start processing transfer on registered wallet %s ...", wallet.Address)
	//		tx, err := service.newIncomingTrasaction(block, btcTx, wallet.Address, outAddresses[wallet.Address])
	//		if err != nil {
	//			log.Printf("->   Error: creating incoming tx %s for wallet %s: %s", btcTx.Txid, wallet.Address, err)
	//			continue
	//		}
	//
	//		oldTx, err := service.dbConn.FindFirstIncomingTx(tx.Currency, tx.TransactionId, tx.AddressString)
	//		if err != nil {
	//			log.Printf("->   Error: searching incoming tx %s for wallet %s: %s", btcTx.Txid, wallet.Address, err)
	//			continue
	//		}
	//
	//		if oldTx != nil {
	//			// tx.Id = oldTx.Id // transaction will be updated in PutIncomingTx...
	//			// but it was already been converted. So updation must be prevented:
	//			log.Printf("->   Error: tx %s for %s wallet %s HAS ALREADY BEEN ADDED", btcTx.Txid, service.conf.Ticker, wallet.Address)
	//			continue
	//		}
	//
	//		err = service.dbConn.PutIncomingTx(*tx)
	//		if err != nil {
	//			log.Printf("->   Error: inserting incoming tx %s for wallet %s: %s", btcTx.Txid, wallet.Address, err)
	//			continue
	//		}
	//
	//		log.Printf("->   Transfer %g %s to %s has been registered!", outAddresses[wallet.Address], service.conf.Ticker, wallet.Address)
	//	}
	//
	}

	//log.Printf("Block %d has been finished", block.Height)
	//
	//chainState, err := service.dbConn.GetChainState(service.conf.Ticker)
	//if err != nil {
	//	log.Printf("Getting chainstate for %s error: %s", service.conf.Ticker, err)
	//	return
	//}
	//
	//if chainState == nil {
	//	chainState = new(models.ChainState)
	//	chainState.ChainType = service.conf.Ticker
	//}
	//
	//chainState.LastBlock = int64(block.Height) + 1
	//chainState.Timestamp = time.Now()
	//
	//err = service.dbConn.PutChainState(*chainState)
	//if err != nil {
	//	log.Printf("Updating chainstate for %s error: %s", service.conf.Ticker, err)
	//}
	//
	//service.conf.BTCConfig.StartBlockHeight = chainState.LastBlock

	return
}

func collectToAddressesFromTx(txs []*types.Transaction) (addresses map[string]float64) {
	addresses = map[string]float64{}
	//for tx := range txs {
	//	if len(tx) == 0 {
	//		// log.Printf("tx %s: Skipped zero out-address", tx.Txid)
	//		continue
	//	}
	//	if len(out.ScriptPubKey.Addresses) > 1 {
	//		log.Printf("tx %s: Skipped multi-address out-address:", tx.Txid)
	//		for _, addr := range out.ScriptPubKey.Addresses {
	//			fmt.Printf("\t%s", addr)
	//		}
	//		continue
	//	}
	//
	//	address := out.ScriptPubKey.Addresses[0]
	//	value := out.Value
	//	addresses[address] = value
	//}
	return
}

