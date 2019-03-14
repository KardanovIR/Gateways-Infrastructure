package services

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func TestNodeClient_GetTransactionByTx(t *testing.T) {
	ctx, log := beforeTest()
	// start test

	// check fee and transfered amount
	tx, err := GetNodeClient().GetTransactionByTx(ctx, "BzDFtaQdnWpWpQsSbPZfzTLqKazHWkuAF7rYCRU8P3Wq")
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
