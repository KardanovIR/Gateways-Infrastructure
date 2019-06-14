package services

import (
	"testing"
)

func TestNodeClient_ValidateAddress(t *testing.T) {
	ctx, log := beforeTest()
	ok, err := GetNodeClient().ValidateAddress(ctx, "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs")
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
	ok, err = GetNodeClient().ValidateAddress(ctx, "9i8x9d4KUVF3Xs9Ks8hmAmLrjoZy6Df2xW3kBnTGkDPCk9yBXmt")
	if err == nil || ok {
		t.Fail()
	}
}
