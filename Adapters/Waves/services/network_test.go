package services

import (
	"testing"
)

func TestNodeClient_Fee(t *testing.T) {
	ctx, log := beforeTest()
	// fee for waves
	fee, err := GetNodeClient().Fee(ctx, "4eWSUDjoYnp2Y4J6vTbY1LY2wT9BZAznTvsv6PM14iaM", "")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if fee != 100000 {
		t.Fail()
	}

	// wBTC is not sponsored -> must fails
	_, err = GetNodeClient().Fee(ctx, "4eWSUDjoYnp2Y4J6vTbY1LY2wT9BZAznTvsv6PM14iaM", "B47PzFMea7HUjpZ8BYwWhPTkCNrdXutAhX8L9Z9tSBdq")
	if err == nil {
		t.FailNow()
	}

}
