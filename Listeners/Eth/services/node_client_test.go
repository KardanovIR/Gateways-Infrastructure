package services

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/config"
)

// add NODE_HOST parameter to env variable (node should be parity)
func TestNodeClient_GetEthRecipientsForTxIncludeInternal(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	err := config.Load("./../listener_test/testdata/config_test.yml")
	if err != nil {
		log.Fatal(err)
	}
	nodeClient, err := newNodeClient(ctx, config.Cfg.Node.Host)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	tx, _, _ := nodeClient.TransactionByHash(ctx, common.HexToHash("0xfa4c3cf5cf4578a5b051039db0b20061471fafd09bfaa388f93b74b79f03f372"))
	data := hexutil.Encode(tx.Data())
	fmt.Println(data)

	//erc-20 transfer
	r, err := nodeClient.GetEthRecipientsForTxIncludeInternal(ctx, "0x380159a8091c991c24b9539d3b42914e6b53a612ffa96ecbec450a973a36b322")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, 0, len(r))
	// normal transaction
	r2, err := nodeClient.GetEthRecipientsForTxIncludeInternal(ctx, "0x6d085f77c7f2eb29062ebd9df13695dff8d5800b5f905c0be2fa930ad7d55fc3")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, 1, len(r2))
	assert.Equal(t, "0x24C2475A0A9d25f393fa8060A39A931D91d5293A", r2[0])

	// internal tx
	r3, err := nodeClient.GetEthRecipientsForTxIncludeInternal(ctx, "0xcc4dd1b4f9e437ef30bdc535e5997115f88948ce405a80b654de73ec6169693e")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, 1, len(r3))
	assert.Equal(t, "0xd7D245547252eF44bf551325C28291848b6af90F", r3[0])
}
