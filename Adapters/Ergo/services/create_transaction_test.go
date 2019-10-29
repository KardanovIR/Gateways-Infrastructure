package services

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/models"
)

func TestCreateInputs_addSmallInputsToTx(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	var unspentInputs = make([]models.UnSpentTx, 0)
	file, err := ioutil.ReadFile("./testdata/unspentInputs.json")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if err := json.Unmarshal(file, &unspentInputs); err != nil {
		log.Error(err)
		t.FailNow()
	}
	unspentInputsForTx := addSmallInputsToTx(findInputs(unspentInputs, 100000))
	assert.Equal(t, uint64(100000), unspentInputsForTx[0].Value)
	assert.Equal(t, "6", unspentInputsForTx[0].ID)
	for i := 1; i < 22; i++ {
		assert.Equal(t, uint64(31000), unspentInputsForTx[i].Value)
	}
	assert.Equal(t, uint64(32000), unspentInputsForTx[23].Value)
	assert.Equal(t, uint64(34000), unspentInputsForTx[24].Value)
	assert.Equal(t, uint64(35000), unspentInputsForTx[25].Value)
	assert.Equal(t, uint64(38000), unspentInputsForTx[26].Value)
	assert.Equal(t, uint64(50000), unspentInputsForTx[27].Value)
	assert.Equal(t, uint64(50250), unspentInputsForTx[28].Value)
	assert.Equal(t, uint64(60000), unspentInputsForTx[29].Value)
}

func TestCreateInputs_createTxInputs(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	var unspentInputs = make([]models.UnSpentTx, 0)

	file, err := ioutil.ReadFile("./testdata/unspentInputs.json")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if err := json.Unmarshal(file, &unspentInputs); err != nil {
		log.Error(err)
		t.FailNow()
	}
	txInputs, summaryAmount, err := createTxInputs(ctx, unspentInputs, 100000)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1081250), summaryAmount)
	assert.Equal(t, 30, len(txInputs))
}

func TestCreateInputs_addSmallInputsToTx2(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	var unspentInputs = make([]models.UnSpentTx, 0)
	file, err := ioutil.ReadFile("./testdata/unspentInputs.json")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if err := json.Unmarshal(file, &unspentInputs); err != nil {
		log.Error(err)
		t.FailNow()
	}
	unspentInputsForTx2 := addSmallInputsToTx(findInputs(unspentInputs, 1693251))
	assert.Equal(t, 30, len(unspentInputsForTx2))
}

func TestCreateInputs_InsufficientFunds(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	var unspentInputs = make([]models.UnSpentTx, 0)
	file, err := ioutil.ReadFile("./testdata/unspentInputs.json")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if err := json.Unmarshal(file, &unspentInputs); err != nil {
		log.Error(err)
		t.FailNow()
	}
	_, _, err = createTxInputs(ctx, unspentInputs, 1693200)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "insufficient funds"))

}

func TestCreateInputs(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	var unspentInputs = make([]models.UnSpentTx, 0)
	file, err := ioutil.ReadFile("./testdata/unspentInputs2.json")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if err := json.Unmarshal(file, &unspentInputs); err != nil {
		log.Error(err)
		t.FailNow()
	}
	txInputs, summaryAmount, err := createTxInputs(ctx, unspentInputs, 4000000)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(txInputs))
	assert.Equal(t, uint64(484679910), summaryAmount)

}
