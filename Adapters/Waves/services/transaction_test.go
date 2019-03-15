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
	_, err := GetNodeClient().CreateRawTxBySendersPublicKey(ctx, "7XM5z1CrfRP6byT5GLPdqQADc35HQ8u6PBE4rXPBB2z5", "3PHgJvRr5EinAB2hVFAPF83xeNZvp7owY9W", 10000000)
	if err == nil {
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
