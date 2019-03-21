package services

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/models"
)

const (
	privateKey = "f8a425216e5b765a3abc829eeb1c0fb8fe291fe9612b23818ee201a7c7a276e8"
	address    = "0xfFDB407fD780b62f43303cCC1f8B0ecF46c72e5b"
	// ChainLink Token contract
	contractAddress = "0x20fE562d797A42Dcb3399062AE9546cd06f63280"
)

var initTest sync.Once

// test for following process:
// 1) generate receiver address
// 2) send money from predefined address to generated address
// 3) return money back: send money from generated address to predefined address
//
func TestNodeClient_Transactions(t *testing.T) {
	ctx, log := beforeTest()
	amount, _ := new(big.Int).SetString("100000000000000", 10)
	// check fee and transfered amount
	fee, err := GetNodeClient().SuggestFee(ctx)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if fee.Cmp(amount) >= 0 {
		log.Error("fee %s more than sending amount %s", fee, amount)
		t.FailNow()
	}
	// check sender's balance
	b, err := GetNodeClient().GetEthBalance(ctx, address)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	amountPlusFee := new(big.Int).Add(amount, fee)
	if b.Cmp(amountPlusFee) <= 0 {
		log.Error("balance %s on sender's address is not more than sending amount %s plus fee %s", b, amount, fee)
		t.FailNow()
	}
	// generate receiver address
	address2, err := GetNodeClient().GenerateAddress(ctx)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.True(t, len(cl.(*nodeClient).privateKeys) > 0)
	log.Infof("Private hex %s, public address %s",
		hex.EncodeToString(crypto.FromECDSA(cl.(*nodeClient).privateKeys[address2])), address2)

	// send 0.000001 ETH to receiver
	tx, err := GetNodeClient().CreateRawTransaction(ctx, address, address2, amount)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	signedTx, err := GetNodeClient().SignTransaction(ctx, address, tx)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	txId, err := GetNodeClient().SendTransaction(ctx, signedTx)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	log.Infof("send transaction %s", txId)
	// wait while transaction will be complete
	err = waitForTxComplete(ctx, txId)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	// check receiver's balance
	balance, err := GetNodeClient().GetEthBalance(ctx, address2)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, balance, amount)

	// return money back
	fee2, err := GetNodeClient().SuggestFee(ctx)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	amountBack := new(big.Int).Sub(balance, fee2)
	tx2, err := GetNodeClient().CreateRawTransaction(ctx, address2, address, amountBack)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	signedTx2, err := GetNodeClient().SignTransaction(ctx, address2, tx2)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	txId2, err := GetNodeClient().SendTransaction(ctx, signedTx2)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	// wait while transaction will be complete
	err = waitForTxComplete(ctx, txId2)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
}

func TestNodeClient_SendErc20(t *testing.T) {
	ctx, log := beforeTest()
	// send 0.0001 Link to receiver
	amount, _ := new(big.Int).SetString("100000000000000", 10)
	address1 := "0x5F862eff5Fb0F2b6B3d83F714A8fe581a8d78e62"
	privateKey1 := "63e8892145d22fbff2c8381b242bd78f191b32c718d51b99dd6b7f9319822320"

	address2 := "0x7D7EB567Df197471A3C43e504844883538356635"
	//privateKey2 := "f0229190763eb29915c40f1e439f510461ec31d6228eceb434fad13659aef0c1"

	tx, err := GetNodeClient().CreateErc20TokensRawTransaction(ctx, address1, contractAddress, address2, amount)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	signedTx, err := GetNodeClient().SignTransactionWithPrivateKey(ctx, privateKey1, tx)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}

	txId, err := GetNodeClient().SendTransaction(ctx, signedTx)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	err = waitForTxComplete(ctx, txId)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
}

func beforeTest() (context.Context, logger.ILogger) {
	ctx := context.Background()
	initTest.Do(func() {
		log, _ := logger.Init(false, logger.DEBUG)
		err := config.Load("./testdata/config_test.yml")
		if err != nil {
			log.Fatal(err)
		}
		err = New(ctx, config.Cfg.Node)
		if err != nil {
			log.Fatal(err)
		}
		pk, err := crypto.HexToECDSA(privateKey)
		if err != nil {
			log.Fatal("can't cast key to ECDSA: %s", err)
		}
		cl.(*nodeClient).privateKeys[address] = pk
	})
	log := logger.FromContext(ctx)
	return ctx, log
}

func waitForTxComplete(ctx context.Context, txId string) error {
	log := logger.FromContext(ctx)
	// wait while transaction will be complete
	for i := 0; i < 20; i++ {
		lastStatus, err := GetNodeClient().GetTxStatusByTxID(ctx, txId)
		if err != nil {
			log.Error(err)
			return err
		}
		if lastStatus == models.TxStatusUnKnown {
			log.Error("unknown transaction by txID %s ", txId)
			return err
		}
		if lastStatus == models.TxStatusSuccess {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return errors.New("transaction in pending status yet")
}
