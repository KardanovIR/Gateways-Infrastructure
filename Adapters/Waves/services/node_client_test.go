package services

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
)

func TestGetLastBlockHeight(t *testing.T) {
	ctx, _ := beforeTest()
	// setup

	block, err := GetNodeClient().GetLastBlockHeight(ctx)
	if err != nil {
		t.Fail()
	}
	bl, err := strconv.Atoi(block)
	if bl < 533148 {
		t.Fail()
	}
}

const (
	privateKey = "AAA9yc4jsbN8hTGCzygxkoKbCYwgs7SuqAwbU6cb1nhi"
	publicKey  = "7XM5z1CrfRP6byT5GLPdqQADc35HQ8u6PBE4rXPBB2z5"
	address    = "3N1cBHN9L3YFuYuJXXpnFVu67Vw726wPZ5Y"
)

const pause = time.Second * 10

func TestNodeClient(t *testing.T) {
	ctx, log := beforeTest()
	// start test
	amount := uint64(1000000)
	// check fee and transfered amount
	fee, err := GetNodeClient().Fee(ctx, publicKey, "")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if fee > amount {
		log.Errorf("fee %s more than sending amount %s", fee, amount)
		t.FailNow()
	}
	// check sender's balance
	balance, err := GetNodeClient().GetBalance(ctx, address)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	amountPlusFee := amount + fee
	if balance < amountPlusFee {
		log.Errorf("balance %d on sender's address is not more than sending amount %d plus fee %d", balance, amount, fee)
		t.FailNow()
	}
	// generate receiver address
	address2, err := GetNodeClient().GenerateAddress(ctx)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	assert.True(t, len(cl.(*nodeClient).privateKeys) > 0)
	log.Infof("Private hex %s, address %s", cl.(*nodeClient).privateKeys[address2], address2)

	// send 0.001 WAVES to receiver
	tx, err := GetNodeClient().CreateRawTxBySendersPublicKey(ctx, publicKey, address2, amount, "")
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	signedTx, _, err := GetNodeClient().SignTxWithSecretKey(ctx, privateKey, tx)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	txId, err := GetNodeClient().SendTransaction(ctx, signedTx)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	log.Infof("send transaction %s", txId)
	// wait while transaction will be complete
	var i = 0
	for ; i < 10; i++ {
		newBalance, err := GetNodeClient().GetBalance(ctx, address2)
		log.Infof("transaction balance %d", newBalance)
		if err != nil {
			log.Error(err)
			t.Fail()
		}

		if newBalance == amount {
			break
		}
		time.Sleep(pause)
	}
	if i == 10 {
		log.Error("transaction doesn't complete!! %s ", txId)
		t.FailNow()
	}

	// return money back
	amountBack := amount - fee
	tx2, err := GetNodeClient().CreateRawTxBySendersAddress(ctx, address2, address, amountBack)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	signedTx2, _, err := GetNodeClient().SignTxWithKeepedSecretKey(ctx, address2, tx2)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	txId2, err := GetNodeClient().SendTransaction(ctx, signedTx2)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	log.Infof("send transaction %s", txId2)
	// wait while transaction will be complete
	var newBalance uint64
	for i := 0; i < 10; i++ {
		newBalance, err = GetNodeClient().GetBalance(ctx, address)
		log.Infof("transaction balance %d", newBalance)
		if err != nil {
			log.Error(err)
			t.Fail()
		}

		if newBalance == balance-2*fee {
			break
		}
		time.Sleep(pause)
	}
	assert.Equal(t, newBalance, balance-2*fee)
}

func beforeTest() (context.Context, logger.ILogger) {
	ctx := context.Background()
	log, _ := logger.Init(false, logger.DEBUG)
	err := config.Load("./testdata/config_test.yml")
	if err != nil {
		log.Fatal(err)
	}
	if err := New(ctx, config.Cfg.Node); err != nil {
		log.Fatal("can't create node's client: ", err)
	}
	return ctx, log
}
