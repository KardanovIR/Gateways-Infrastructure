package services

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
)

const txTraceMethodName = "trace_transaction"

var bigIntZero = big.NewInt(0)

type nodeClient struct {
	*ethclient.Client
	rpcNodeClient *rpc.Client
}

func newNodeClient(ctx context.Context, nodeHost string) (*nodeClient, error) {
	log := logger.FromContext(ctx)
	rpcc, err := newRPCClient(log, nodeHost)
	if err != nil {
		log.Errorf("error during initialise rpc client: %s", err)

		return nil, err
	}
	ethClient := ethclient.NewClient(rpcc)
	return &nodeClient{ethClient, rpcc}, nil
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

// GetEthRecipientsForTxIncludeInternal parse eth transfers include internal transactions. Work only with parity node
func (n *nodeClient) GetEthRecipientsForTxIncludeInternal(ctx context.Context, txHash string) ([]string, error) {
	log := logger.FromContext(ctx)
	traceList, err := n.getTxTrace(ctx, txHash)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	recipients := n.findRecipientsAddresses(ctx, traceList)
	return recipients, nil
}

type Trace struct {
	Action Action `json:"action"`
}

type Action struct {
	CallType string `json:"callType"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
}

func (n *nodeClient) getTxTrace(ctx context.Context, txHash string) ([]Trace, error) {
	log := logger.FromContext(ctx)
	result := make([]Trace, 0)
	if err := n.rpcNodeClient.CallContext(ctx, &result, txTraceMethodName, txHash); err != nil {
		log.Errorf("'trace_transaction' call finished with error: %s", err)
		return result, err
	}
	log.Debugf("eth transfers count in tx %s is %d", txHash, len(result))
	return result, nil
}

func (n *nodeClient) findRecipientsAddresses(ctx context.Context, traceList []Trace) []string {
	log := logger.FromContext(ctx)
	recipients := make([]string, 0)
	hasAddressFunc := func(searchFor string) bool {
		for _, a := range recipients {
			if a == searchFor {
				return true
			}
		}
		return false
	}
	for _, tr := range traceList {
		amount, ok := new(big.Int).SetString(tr.Action.Value, 0)
		// only for transfers with eth amount
		if ok && amount.Cmp(bigIntZero) > 0 {
			// convert address string -> object -> string to get address with right letters case
			addr := common.HexToAddress(tr.Action.To)
			if !hasAddressFunc(addr.Hex()) {
				recipients = append(recipients, addr.Hex())
			}
		}
		if !ok {
			log.Errorf("can't convert %s to big Int", tr.Action.Value)
		}
	}
	log.Debugf("not zero eth transfers count is %d", len(recipients))
	return recipients
}
