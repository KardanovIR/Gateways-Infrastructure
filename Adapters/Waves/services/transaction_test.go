package services

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/models"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func TestNodeClient_GetTransactionByTxId(t *testing.T) {
	ctx, log := beforeTest()
	tx, err := GetNodeClient().GetTransactionByTxId(ctx, "BzDFtaQdnWpWpQsSbPZfzTLqKazHWkuAF7rYCRU8P3Wq")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	tr := new(proto.TransferV2)
	if err := tr.UnmarshalBinary(tx); err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, tr.ID.String(), "BzDFtaQdnWpWpQsSbPZfzTLqKazHWkuAF7rYCRU8P3Wq")
	assert.Equal(t, tr.Recipient.String(), "3NAQRd3SnKBqnUhf75SjBGybiJBm9kVyYYJ")
	assert.Equal(t, tr.Amount, uint64(1000000))
}

func TestNodeClient_CreateRawTxBySendersPublicKey(t *testing.T) {
	ctx, log := beforeTest()
	// send to incorrect address -> must fails
	_, err := GetNodeClient().CreateRawTxBySendersPublicKey(ctx, "7XM5z1CrfRP6byT5GLPdqQADc35HQ8u6PBE4rXPBB2z5",
		"3PHgJvRr5EinAB2hVFAPF83xeNZvp7owY9W", 10000000, "")
	if err == nil {
		log.Error(err)
		t.Fail()
	}
}

// send wBTC from 1 to 2, if money on 1 account is finished - return them from 2 account
func TestNodeClient_SendWBtc(t *testing.T) {
	ctx, log := beforeTest()
	// send asset bitcoin from 1 account to 2 account
	bitcoinAssetId := "DWgwcZTMhSvnyYCoWLRUXXSH1RSkzThXLJhww9gwkqdn"
	var (
		privateKey1 = "DyjaYCk9U1CMNbUgcq64WHeasnqAxw58gBkxwKmZJ3ob"
		publicKey1  = "3TQRApHb85CR8A1eWKgRdkgLqbkUCb3BFXQhw8bx79Wb"
	//	address1    = "3N378WXBCUFVesBsqof5ra9EHbvwJPYPtYM"
	)
	var (
		//	privateKey2 = "DwCHooqdi9rsbLy87bvT8ywuNvZbmvJURYgGJqGjaDKH"
		//	publicKey2  = "Bf5DKiZCWHBdWkjTUhrzQ9qHw7cxmTzb7dDMTvf8Zw6z"
		address2 = "3N1eTbkh9RikjekLoqisfhKw7gsAquPZAvM"
	)
	tx, err := GetNodeClient().CreateRawTxBySendersPublicKey(
		ctx, publicKey1, address2, 100, bitcoinAssetId)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	signedTx, err := GetNodeClient().SignTxWithSecretKey(ctx, privateKey1, tx)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	_, err = GetNodeClient().SendTransaction(ctx, signedTx)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
}

func TestNodeClient_GetTransactionStatus(t *testing.T) {
	ctx, log := beforeTest()
	// send to incorrect address -> must fails
	st, err := GetNodeClient().GetTransactionStatus(ctx, "4xjtyGErZ528b39aNrhBqoTq9RYQnZreEaP4QfbTMFjG")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, st, models.TxStatusSuccess)
	st2, err := GetNodeClient().GetTransactionStatus(ctx, "24XscwjpC113ijAzQkTaUJFpumQA7XWtUVydBhBxtWMv")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, st2, models.TxStatusUnKnown)
}
