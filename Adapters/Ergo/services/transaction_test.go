package services

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Ergo/services/converter"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
)

func TestNodeClient_replaceQuotesFromSides(t *testing.T) {
	s := "\"c99ce66bd0cebf45f97aba2f48912583562a7cc8fdf4c89079608517b1955c73\""
	assert.Equal(t, "c99ce66bd0cebf45f97aba2f48912583562a7cc8fdf4c89079608517b1955c73", replaceQuotesFromSides(s))
	s2 := "c99ce66bd0cebf45f97aba2f48912583562a7cc8fdf4c89079608517b1955c73"
	assert.Equal(t, "c99ce66bd0cebf45f97aba2f48912583562a7cc8fdf4c89079608517b1955c73", replaceQuotesFromSides(s2))
}

func TestNodeClient_parseTx(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	converter.Init(ctx, 8)
	tx := models.Tx{}
	//parseTx_1.json - 2 inputs with same addresses with charge

	file1, err := ioutil.ReadFile("./testdata/parseTx_1.json")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if err := json.Unmarshal(file1, &tx); err != nil {
		log.Error(err)
		t.FailNow()
	}
	txInfo := parseTx(&tx)
	assert.Equal(t, 1, len(txInfo.Inputs))
	assert.Equal(t, 1, len(txInfo.Outputs))
	assert.Equal(t, "130000", txInfo.Outputs[0].Amount)
	assert.Equal(t, "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs", txInfo.Outputs[0].Address)

	assert.Equal(t, "230000", txInfo.Inputs[0].Amount)
	assert.Equal(t, "3WwgqLMBZhUWVHQUoYakSmcJwte8TPYM3gFkYeJ84S3NP21T2uJg", txInfo.Inputs[0].Address)

	assert.Equal(t, "100000", txInfo.Fee)
	assert.Equal(t, "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs", txInfo.To)
	assert.Equal(t, "3WwgqLMBZhUWVHQUoYakSmcJwte8TPYM3gFkYeJ84S3NP21T2uJg", txInfo.From)
	assert.Equal(t, "bbd9118ce330aeea59f9768bc0d1b9d0e73fec63abc7c60ed7b5f444ba16c8e0", txInfo.TxHash)
	assert.Equal(t, "130000", txInfo.Amount)

	// case2
	//parseTx_2.json - 2 inputs with same addresses without charge
	tx2 := models.Tx{}
	file2, err := ioutil.ReadFile("./testdata/parseTx_2.json")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if err := json.Unmarshal(file2, &tx2); err != nil {
		log.Error(err)
		t.FailNow()
	}
	txInfo2 := parseTx(&tx2)
	assert.Equal(t, 1, len(txInfo2.Inputs))
	assert.Equal(t, 1, len(txInfo2.Outputs))
	assert.Equal(t, "140000", txInfo2.Outputs[0].Amount)
	assert.Equal(t, "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs", txInfo2.Outputs[0].Address)

	assert.Equal(t, "240000", txInfo2.Inputs[0].Amount)
	assert.Equal(t, "3WwgqLMBZhUWVHQUoYakSmcJwte8TPYM3gFkYeJ84S3NP21T2uJg", txInfo2.Inputs[0].Address)

	assert.Equal(t, "100000", txInfo2.Fee)
	assert.Equal(t, "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs", txInfo2.To)
	assert.Equal(t, "3WwgqLMBZhUWVHQUoYakSmcJwte8TPYM3gFkYeJ84S3NP21T2uJg", txInfo2.From)
	assert.Equal(t, "cd52cac69416be9205096d02a1fe59ef17aae6afb03fbec8ff91406e95b1318d", txInfo2.TxHash)
	assert.Equal(t, "140000", txInfo2.Amount)

	//parseTx_3.json - 1 input with charge
	tx3 := models.Tx{}
	file3, err := ioutil.ReadFile("./testdata/parseTx_3.json")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if err := json.Unmarshal(file3, &tx3); err != nil {
		log.Error(err)
		t.FailNow()
	}
	txInfo3 := parseTx(&tx3)
	assert.Equal(t, 1, len(txInfo3.Inputs))
	assert.Equal(t, 1, len(txInfo3.Outputs))
	assert.Equal(t, "120000", txInfo3.Outputs[0].Amount)
	assert.Equal(t, "3WwgqLMBZhUWVHQUoYakSmcJwte8TPYM3gFkYeJ84S3NP21T2uJg", txInfo3.Outputs[0].Address)

	assert.Equal(t, "220000", txInfo3.Inputs[0].Amount)
	assert.Equal(t, "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs", txInfo3.Inputs[0].Address)

	assert.Equal(t, "100000", txInfo3.Fee)
	assert.Equal(t, "3WwgqLMBZhUWVHQUoYakSmcJwte8TPYM3gFkYeJ84S3NP21T2uJg", txInfo3.To)
	assert.Equal(t, "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs", txInfo3.From)
	assert.Equal(t, "1857c4e2490ff80cec9dc2ffdf64fb367744130c39641106562f88cf696f5096", txInfo3.TxHash)
	assert.Equal(t, "120000", txInfo3.Amount)

	//parseTx_4.json - 2 different address in input with charge on 1 input and output to 1 input
	tx4 := models.Tx{}
	file4, err := ioutil.ReadFile("./testdata/parseTx_4.json")
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	if err := json.Unmarshal(file4, &tx4); err != nil {
		log.Error(err)
		t.FailNow()
	}
	txInfo4 := parseTx(&tx4)
	assert.Equal(t, 1, len(txInfo4.Inputs))
	assert.Equal(t, 2, len(txInfo4.Outputs))
	assert.Equal(t, "25", txInfo4.Inputs[0].Amount)
	assert.Equal(t, "3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN", txInfo4.Inputs[0].Address)
	assert.Equal(t, "1", txInfo4.Fee)
	if txInfo4.Outputs[0].Address == "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs" {
		checkFirstOutPut(t, txInfo4.Outputs[0])
		checkSecondOutPut(t, txInfo4.Outputs[1])
	} else {
		checkFirstOutPut(t, txInfo4.Outputs[1])
		checkSecondOutPut(t, txInfo4.Outputs[0])
	}
	assert.Equal(t, "", txInfo4.To)
	assert.Equal(t, "3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN", txInfo4.From)
	assert.Equal(t, "bbd9118ce330aeea59f9768bc0d1b9d0e73fec63abc7c60ed7b5f444ba16c8e0", txInfo4.TxHash)
	assert.Equal(t, "", txInfo4.Amount)
}

func checkFirstOutPut(t *testing.T, info models.InputOutputInfo) {
	assert.Equal(t, "21", info.Amount)
	assert.Equal(t, "3WwHhExDYkWrkjpqe3BuH4FSAzMeMkxZiuhwRpNUoBJrD7BbJpzs", info.Address)
}

func checkSecondOutPut(t *testing.T, info models.InputOutputInfo) {
	assert.Equal(t, "3", info.Amount)
	assert.Equal(t, "3WwgqLMBZhUWVHQUoYakSmcJwte8TPYM3gFkYeJ84S3NP21T2uJg", info.Address)
}
