package services

import (
	"context"
	"crypto/ecdsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeClient_IsAddressValid(t *testing.T) {
	cl := nodeClient{}
	assert.True(t, cl.IsAddressValid(context.Background(), "0x74d2d6195a1c374e8043920bf7530f7750ec3c5d"))
	assert.False(t, cl.IsAddressValid(context.Background(), "1Po1oWkD2LmodfkBYiAktwh76vkF93LKnh"))
	assert.False(t, cl.IsAddressValid(context.Background(), "2N3sWVq5inguiqmyzZpSQKfXqwtWTDnre7p"))
}

func TestNodeClient_GenerateAddress(t *testing.T) {
	cl := nodeClient{chainID: 3, privateKeys: make(map[string]*ecdsa.PrivateKey)}
	pb, err := cl.GenerateAddress(context.Background())
	if err != nil {
		t.Fail()
	}
	assert.True(t, cl.IsAddressValid(context.Background(), pb))
}
