package adapter_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/clientgrpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/services"
)

const (
	privateKey = "2eN6rsKcTnykyPsppwHnNmB4WQNc86ZwbJto3sXPZCYf"
	publicKey  = "3FdG1P3KzLxNgGW9BxLGpzk9G8rKWDujBST8LKjmMMdv"
	address    = "3MwJ9fszo2VNgD5kuXMs5nypcsCPSukmSBA"
)

// TestGrpcClient checks all endpoints of Waves adapter sending requests via grpc client
// 1) Fee method
// 2) GetBalance of predefined account
// 3) check that money on this address is enough for test
// 4) GenerateAddress method: getting receiver address
// 5) ValidateAddress: check address generated on previous step
// 6) GetRawTransactionBySendersPublicKey: get transaction to send money from predefined address to generated address
// 7) SignTransactionBySecretKey: sign transaction created on previous step
// 8) SendTransaction: send transaction signed on previous step
// 9)  GetTransactionStatus: check status by hash get on previous step (do it on loop with 10 seconds pause)
// 10) GetBalance of receiver account: it must be equal transfer amount
//     and send transaction from generated address to predefined (return money back):
//     transfered amount = balance_on_generated_address - tFee()
//     GetRawTransactionBySendersAddress, SignTransaction and send transaction, wait for it's completion
// 11) check balance on predefined account: balance_on_2_step - fee_on_1_step - fee_on_10_step
func TestGrpcClient(t *testing.T) {
	ctx, log := beforeTests()
	amount, _ := new(big.Int).SetString("1000000", 10)

	// check fee and transfered amount on predefined address
	feeReply, err := clientgrpc.GetClient().Fee(ctx, &wavesAdapter.FeeRequest{
		SendersPublicKey: "3FdG1P3KzLxNgGW9BxLGpzk9G8rKWDujBST8LKjmMMdv", AssetId: ""})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	fee, _ := new(big.Int).SetString(feeReply.Fee, 10)
	if fee.Cmp(amount) >= 0 {
		log.Errorf("fee %s more than sending amount %s", fee, amount)
		t.FailNow()
	}
	// check sender's balance
	b, err := clientgrpc.GetClient().GetBalance(ctx, &wavesAdapter.AddressRequest{Address: address})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	sBalance, _ := new(big.Int).SetString(b.Balance, 10)
	amountPlusFee := new(big.Int).Add(amount, fee)
	if sBalance.Cmp(amountPlusFee) <= 0 {
		log.Errorf("balance %s on sender's address is not more than sending amount %s plus feeReply %s", sBalance, amount, fee)
		t.FailNow()
	}
	// generate receiver address
	address2Reply, err := clientgrpc.GetClient().GenerateAddress(ctx, &wavesAdapter.EmptyRequest{})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}

	// check generated address
	isValidReply, err := clientgrpc.GetClient().ValidateAddress(ctx, &wavesAdapter.AddressRequest{Address: address2Reply.Address})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if !assert.True(t, isValidReply.Valid) {
		t.FailNow()
	}
	// send 0.0001 Waves to receiver
	address2 := address2Reply.Address
	tx, err := clientgrpc.GetClient().GetRawTransactionBySendersPublicKey(ctx,
		&wavesAdapter.RawTransactionBySendersPublicKeyRequest{
			SendersPublicKey: publicKey,
			AddressTo:        address2,
			Amount:           amount.String(),
		})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	signedTx, err := clientgrpc.GetClient().SignTransactionBySecretKey(ctx,
		&wavesAdapter.SignTransactionBySecretKeyRequest{SenderSecretKey: privateKey, Tx: tx.Tx})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	sendTxReply, err := clientgrpc.GetClient().SendTransaction(ctx, &wavesAdapter.SendTransactionRequest{Tx: signedTx.Tx})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	log.Infof("send transaction %s", sendTxReply.TxId)
	if err := waitForTxComplete(ctx, sendTxReply.TxId); err != nil {
		log.Error(err)
		t.FailNow()
	}

	// check receiver's balance
	balanceReply, err := clientgrpc.GetClient().GetBalance(ctx, &wavesAdapter.AddressRequest{Address: address2})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, balanceReply.Balance, amount.String())

	// return money back

	balance, _ := new(big.Int).SetString(balanceReply.Balance, 10)
	amountBack := new(big.Int).Sub(balance, fee)
	tx2, err := clientgrpc.GetClient().GetRawTransactionBySendersAddress(ctx,
		&wavesAdapter.RawTransactionBySendersAddressRequest{
			AddressFrom: address2,
			AddressTo:   address,
			Amount:      amountBack.String(),
		})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	signedTx2, err := clientgrpc.GetClient().SignTransaction(ctx,
		&wavesAdapter.SignTransactionRequest{SenderAddress: address2, Tx: tx2.Tx})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	txId2, err := clientgrpc.GetClient().SendTransaction(ctx, &wavesAdapter.SendTransactionRequest{Tx: signedTx2.Tx})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	log.Infof("send transaction %s", txId2)
	// wait while transaction will be complete
	if err := waitForTxComplete(ctx, txId2.TxId); err != nil {
		log.Error(err)
		t.FailNow()
	}

	// check balance
	balance1Reply, err := clientgrpc.GetClient().GetBalance(ctx, &wavesAdapter.AddressRequest{Address: address})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	// final_balance_on_address1 = initial_balance_on_address1 - (fee_of_1_to_2 + fee_of_2_to_1)
	amountResult := new(big.Int).Sub(sBalance, new(big.Int).Add(fee, fee))
	assert.Equal(t, amountResult.String(), balance1Reply.Balance)
}

func waitForTxComplete(ctx context.Context, txID string) error {
	log := logger.FromContext(ctx)
	// wait while transaction will be complete
	for i := 0; i < 10; i++ {
		statusReply, err := clientgrpc.GetClient().GetTransactionStatus(ctx, &wavesAdapter.GetTransactionStatusRequest{TxId: txID})
		log.Infof("transaction status %s", statusReply)
		if err != nil {
			return err
		}
		if statusReply.Status == string(models.TxStatusUnKnown) {
			return errors.New("unknown transaction")
		}
		if statusReply.Status == string(models.TxStatusSuccess) {
			log.Infof("returned from loop on %d iteration", i+1)
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return errors.New("transaction in pending status yet")
}

func TestGetAllBalances(t *testing.T) {
	ctx, log := beforeTests()
	address := "3N7DGmkCmUgMo9jpuUekUhCMpiBRR1Zm51p"
	balanceReply, err := clientgrpc.GetClient().GetAllBalances(ctx, &wavesAdapter.AddressRequest{Address: address})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	wBtcsAssetId := "DWgwcZTMhSvnyYCoWLRUXXSH1RSkzThXLJhww9gwkqdn"
	assert.Equal(t, balanceReply.Amount, "190000000")
	assert.Equal(t, len(balanceReply.AssetBalances), 1)
	assert.Equal(t, balanceReply.AssetBalances[0].AssetId, wBtcsAssetId)
	assert.Equal(t, balanceReply.AssetBalances[0].Amount, "101000")
}

func beforeTests() (context.Context, logger.ILogger) {
	ctx := context.Background()
	log, _ := logger.Init(false, logger.DEBUG)
	err := config.Load("./testdata/config_adapter_test.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = services.New(ctx, config.Cfg.Node)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := server.InitAndStart(ctx, config.Cfg.Port, services.GetNodeClient()); err != nil {
			log.Fatal("Can't start grpc server", err)
		}
	}()

	if err := clientgrpc.New(ctx, ":"+config.Cfg.Port); err != nil {
		log.Fatal("Can't init grpc client", err)
	}
	return ctx, log
}
