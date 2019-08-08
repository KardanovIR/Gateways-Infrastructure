package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeClient_IsAddressValid(t *testing.T) {
	ctx, _ := beforeTest()
	cl := GetNodeClient()
	ok, msg, err := cl.IsAddressValid(ctx, "0x74d2d6195a1c374e8043920bf7530f7750ec3c5d")
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Empty(t, msg)
	ok, msg, err = cl.IsAddressValid(ctx, "1Po1oWkD2LmodfkBYiAktwh76vkF93LKnh")
	assert.False(t, ok)
	ok, msg, err = cl.IsAddressValid(ctx, "2N3sWVq5inguiqmyzZpSQKfXqwtWTDnre7p")
	assert.False(t, ok)
	// smart contract
	ok, msg, err = cl.IsAddressValid(ctx, "0xD5727f9d8C5b9E4472566683F4e562Ef9B47dCE3")
	assert.False(t, ok)
	assert.NotEmpty(t, msg)
}

func TestNodeClient_GenerateAddress(t *testing.T) {
	ctx, _ := beforeTest()
	pb, err := cl.GenerateAddress(ctx)
	if err != nil {
		t.Fail()
	}
	ok, _, _ := cl.IsAddressValid(ctx, pb)
	assert.True(t, ok)
}
