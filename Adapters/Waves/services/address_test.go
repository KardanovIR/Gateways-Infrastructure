package services

import (
	"testing"
)

func TestNodeClient_GenerateAddress(t *testing.T) {
	ctx, log := beforeTest()
	address, err := GetNodeClient().GenerateAddress(ctx)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	ok, err := GetNodeClient().ValidateAddress(ctx, address)
	if err != nil || !ok {
		log.Error(err)
		t.Fail()
	}
}

func TestNodeClient_ValidateAddress(t *testing.T) {
	ctx, log := beforeTest()
	ok, err := GetNodeClient().ValidateAddress(ctx, "3N5fVy6xD7BPXijHEAuPHfg9HV49ywstqeT")
	if err != nil || !ok {
		log.Error(err)
		t.FailNow()
	}
	// bitcoin address must be not valid
	ok, err = GetNodeClient().ValidateAddress(ctx, "3P4dudDfyYiuW7J3JexMYdGMwUUuyUhHQz")
	if err == nil || ok {
		t.Fail()
	}

	// must be not valid because of wrong network
	ok, err = GetNodeClient().ValidateAddress(ctx, "3PHgJvRr5EinAB2hVFAPF83xeNZvp7owY9W")
	if err == nil || ok {
		t.Fail()
	}
}
