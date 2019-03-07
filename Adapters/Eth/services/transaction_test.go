package services

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/magiconair/properties/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/models"
)

const (
	privateKey = "f8a425216e5b765a3abc829eeb1c0fb8fe291fe9612b23818ee201a7c7a276e8"
	address    = "0xfFDB407fD780b62f43303cCC1f8B0ecF46c72e5b"
)

// test for following process:
// 1) generate receiver address
// 2) send money from predefined address to generated address
// 3) return money back: send money from generated address to predefined address
//
func TestNodeClient_Transactions(t *testing.T) {
	ctx := context.Background()
	// setup
	log, _ := logger.Init(false, logger.DEBUG)
	err := config.Load("./testdata/config_test.yml")
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	err = New(ctx, config.Cfg.Node)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	cl.(*nodeClient).privateKeys[address] = pk
	// start test
	amount, _ := new(big.Int).SetString("100000000000000", 10)

	// check fee and transfered amount
	fee, err := GetNodeClient().SuggestFee(ctx)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	if fee.Cmp(amount) >= 0 {
		log.Error("fee %s more than sending amount %s", fee, amount)
		t.Fail()
	}
	// check sender's balance
	b, err := GetNodeClient().GetBalance(ctx, address)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	amountPlusFee := new(big.Int).Add(amount, fee)
	if b.Cmp(amountPlusFee) <= 0 {
		log.Error("balance %s on sender's address is not more than sending amount %s plus fee %s", b, amount, fee)
		t.Fail()
	}
	// generate receiver address
	address2, err := GetNodeClient().GenerateAddress(ctx)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	assert.Equal(t, len(cl.(*nodeClient).privateKeys), 2)
	log.Infof("Private hex %s, public address %s",
		hex.EncodeToString(crypto.FromECDSA(cl.(*nodeClient).privateKeys[address2])), address2)

	// send 0.000001 ETH to receiver
	tx, err := GetNodeClient().CreateRawTransaction(ctx, address, address2, amount)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	signedTx, err := GetNodeClient().SignTransaction(ctx, address, tx)
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
	var lastStatus models.TxStatus
	var i = 0
	for ; i < 10; i++ {
		lastStatus, err := GetNodeClient().GetTxStatusByTxID(ctx, txId)
		log.Infof("transaction status %s", lastStatus)
		if err != nil {
			log.Error(err)
			t.Fail()
		}
		if lastStatus == models.TxStatusUnKnown {
			log.Error("unknown transaction by txID %s ", txId)
			t.Fail()
		}
		if lastStatus == models.TxStatusSuccess {
			break
		}
		time.Sleep(10 * time.Second)
	}
	log.Info("returned from loop on %d iteration", i+1)
	if lastStatus == models.TxStatusPending {
		log.Error("transaction in pending status yet!! %s ", txId)
	}
	// check receiver's balance
	balance, err := GetNodeClient().GetBalance(ctx, address2)
	if err != nil {
		log.Error(err)
	}
	assert.Equal(t, balance, amount)

	// return money back
	fee2, err := GetNodeClient().SuggestFee(ctx)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	amountBack := new(big.Int).Sub(balance, fee2)
	tx2, err := GetNodeClient().CreateRawTransaction(ctx, address2, address, amountBack)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	signedTx2, err := GetNodeClient().SignTransaction(ctx, address2, tx2)
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
	{
		var lastStatus models.TxStatus
		var i = 0
		for ; i < 10; i++ {
			lastStatus, err := GetNodeClient().GetTxStatusByTxID(ctx, txId2)
			if err != nil {
				log.Error(err)
				t.Fail()
			}
			if lastStatus == models.TxStatusUnKnown {
				log.Error("unknown transaction by txID %s ", txId2)
				t.Fail()
			}
			if lastStatus == models.TxStatusSuccess {
				break
			}
			time.Sleep(10 * time.Second)
		}
		log.Infof("returned from loop on %d iteration", i+1)
		if lastStatus == models.TxStatusPending {
			log.Error("back money transaction in pending status yet!! %s ", txId2)
			t.Fail()
		}

	}
}
