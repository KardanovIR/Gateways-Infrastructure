package adapter_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/clientgrpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/config"
	ethAdapter "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/services"
)

const (
	privateKey = "0c69b368f13f340232ced4d3463a05fc55e9d74a9ba4ecf52edb1d33d1de6239"
	address    = "0x3fe9F0886143dd1AE04413854fe8dEBc3B3E0Ab5"
)

// TestGrpcClient checks all endpoints of Eth adapter sending requests via grpc client
// 1) SuggestFee method
// 2) GetBalance of predefined account
// 3) check that money on this address is enough for test
// 4) GenerateAddress method: getting receiver address
// 5) ValidateAddress: check address generated on previous step
// 6) GetRawTransaction: get transaction to send money from predefined address to generated address
// 7) SignTransaction: sign transaction created on previous step
// 8) SendTransaction: send transaction signed on previous step
// 9)  GetTransactionStatus: check status by hash get on previous step (do it on loop with 10 seconds pause)
// 10) GetBalance of receiver account: it must be equal transfer amount
//     and send transaction from generated address to predefined (return money back):
//     transfered amount = balance_on_generated_address - SuggestFee()
//     create, sign and send transaction, wait for it's completion
// 11) check balance on predefined account: balance_on_2_step - fee_on_1_step - fee_on_10_step
// 12) GetNextNonce method: check nonce on generated account: it must be equal 1
func TestGrpcClient(t *testing.T) {
	ctx, log := beforeTests()
	amount, _ := new(big.Int).SetString("100000000000000", 10)

	// check fee and transfered amount on predefined address
	feeReply, err := clientgrpc.GetClient().SuggestFee(ctx, &ethAdapter.EmptyRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fee, _ := new(big.Int).SetString(feeReply.Fee, 10)
	if fee.Cmp(amount) >= 0 {
		log.Fatal("fee %s more than sending amount %s", fee, amount)
	}
	// check sender's balance
	b, err := clientgrpc.GetClient().GetEthBalance(ctx, &ethAdapter.AddressRequest{Address: address})
	if err != nil {
		log.Fatal(err)
	}
	sBalance, _ := new(big.Int).SetString(b.Balance, 10)
	amountPlusFee := new(big.Int).Add(amount, fee)
	if sBalance.Cmp(amountPlusFee) <= 0 {
		log.Fatal("balance %s on sender's address is not more than sending amount %s plus feeReply %s", sBalance, amount, fee)
	}
	// generate receiver address
	address2Reply, err := clientgrpc.GetClient().GenerateAddress(ctx, &ethAdapter.EmptyRequest{})
	if err != nil {
		log.Fatal(err)
	}

	// check generated address
	isValidReply, err := clientgrpc.GetClient().ValidateAddress(ctx, &ethAdapter.AddressRequest{Address: address2Reply.Address})
	if err != nil {
		log.Fatal(err)
	}
	if !assert.True(t, isValidReply.Valid) {
		return
	}
	// send 0.000001 ETH to receiver
	address2 := address2Reply.Address
	tx, err := clientgrpc.GetClient().GetRawTransaction(ctx, &ethAdapter.RawTransactionRequest{
		AddressFrom: address,
		AddressTo:   address2,
		Amount:      amount.String(),
	})
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := clientgrpc.GetClient().SignTransactionWithPrivateKey(ctx,
		&ethAdapter.SignTransactionWithPrivateKeyRequest{PrivateKey: privateKey, Tx: tx.Tx})
	if err != nil {
		log.Fatal(err)
	}
	sendTxReply, err := clientgrpc.GetClient().SendTransaction(ctx, &ethAdapter.SendTransactionRequest{Tx: signedTx.Tx})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("send transaction %s", sendTxReply.TxHash)
	if err := waitForTxComplete(ctx, sendTxReply.TxHash); err != nil {
		log.Fatal(err)
	}

	// check receiver's balance
	balanceReply, err := clientgrpc.GetClient().GetEthBalance(ctx, &ethAdapter.AddressRequest{Address: address2})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, balanceReply.Balance, amount.String())

	// return money back
	fee2Reply, err := clientgrpc.GetClient().SuggestFee(ctx, &ethAdapter.EmptyRequest{})
	if err != nil {
		log.Fatal(err)
	}
	balance, _ := new(big.Int).SetString(balanceReply.Balance, 10)
	fee2, _ := new(big.Int).SetString(fee2Reply.Fee, 10)
	amountBack := new(big.Int).Sub(balance, fee2)

	tx2, err := clientgrpc.GetClient().GetRawTransaction(ctx, &ethAdapter.RawTransactionRequest{
		AddressFrom: address2,
		AddressTo:   address,
		Amount:      amountBack.String(),
	})
	if err != nil {
		log.Fatal(err)
	}
	signedTx2, err := clientgrpc.GetClient().SignTransaction(ctx, &ethAdapter.SignTransactionRequest{SenderAddress: address2, Tx: tx2.Tx})
	if err != nil {
		log.Fatal(err)
	}
	txHash2, err := clientgrpc.GetClient().SendTransaction(ctx, &ethAdapter.SendTransactionRequest{Tx: signedTx2.Tx})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("send transaction %s", txHash2)
	// wait while transaction will be complete
	if err := waitForTxComplete(ctx, txHash2.TxHash); err != nil {
		log.Fatal(err)
	}

	// check balance
	balance1Reply, err := clientgrpc.GetClient().GetEthBalance(ctx, &ethAdapter.AddressRequest{Address: address})
	if err != nil {
		log.Fatal(err)
	}
	// final_balance_on_address1 = initial_balance_on_address1 - (fee_of_1_to_2 + fee_of_2_to_1)
	amountResult := new(big.Int).Sub(sBalance, new(big.Int).Add(fee, fee2))
	assert.Equal(t, amountResult.String(), balance1Reply.Balance)

	nonceReply, err := clientgrpc.GetClient().GetNextNonce(ctx, &ethAdapter.AddressRequest{Address: address2})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, int64(1), nonceReply.Nonce)
}

func TestTokenBalance(t *testing.T) {
	ctx, log := beforeTests()
	address := "0xb8ebd916689e773c6657de537151476A3a8259fc"
	// ERC-20 (LINK)
	contract1 := "0x20fe562d797a42dcb3399062ae9546cd06f63280"
	// Test Standard Token (TST)
	contract2 := "0x722dd3F80BAC40c951b51BdD28Dd19d435762180"
	balances, err := clientgrpc.GetClient().GetAllBalance(ctx, &ethAdapter.GetAllBalanceRequest{
		Address:   address,
		Contracts: []string{contract1, contract2},
	})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, balances.Amount, "4200000000000000")
	assert.Equal(t, len(balances.TokenBalances), 2)
	if balances.TokenBalances[0].Contract == contract1 {
		assert.Equal(t, balances.TokenBalances[0].Amount, "325000000000000000")
		assert.Equal(t, balances.TokenBalances[1].Amount, "19000000000000000000")
		assert.Equal(t, balances.TokenBalances[1].Contract, contract2)
	} else {
		assert.Equal(t, balances.TokenBalances[0].Amount, "19000000000000000000")
		assert.Equal(t, balances.TokenBalances[1].Amount, "325000000000000000")
		assert.Equal(t, balances.TokenBalances[1].Contract, contract1)
	}
}

func waitForTxComplete(ctx context.Context, txHash string) error {
	log := logger.FromContext(ctx)
	// wait while transaction will be complete
	var i = 0
	for ; i < 20; i++ {
		statusReply, err := clientgrpc.GetClient().GetTransactionStatus(ctx, &ethAdapter.GetTransactionStatusRequest{TxHash: txHash})
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
		time.Sleep(10 * time.Second)
	}
	return errors.New("transaction in pending status yet")
}

func beforeTests() (context.Context, logger.ILogger) {
	ctx := context.Background()
	log, _ := logger.Init(false, logger.DEBUG)
	err := config.Load("./testdata/config_test.yml")
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
