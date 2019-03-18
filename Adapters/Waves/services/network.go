package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

const calculateFeeUrl = "transactions/calculateFee"

type FeeResponse struct {
	FeeAssetId *string `json:"feeAssetId"`
	FeeAmount  uint64  `json:"feeAmount"`
}

func (cl *nodeClient) GetLastBlockHeight(ctx context.Context) (string, error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'GetLastBlockHeight'")

	lastBlock, _, err := cl.nodeClient.Blocks.Last(ctx)
	if err != nil {
		log.Errorf("get last block fails: %s", err)
		return "", err
	}
	return strconv.FormatUint(lastBlock.Height, 10), nil
}

func (cl *nodeClient) Fee(ctx context.Context, senderPublicKey string, assetId string) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Debugf("call service method 'Fee' for sender %s, assetId %s", senderPublicKey, assetId)
	senderPublic, err := crypto.NewPublicKeyFromBase58(senderPublicKey)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	address, err := proto.NewAddressFromPublicKey(cl.chainID.Schema(), senderPublic)
	if err != nil {
		return 0, err
	}
	tx, err := createRawTransactionWithoutFee(ctx, senderPublic, address.String(), 1, "", assetId)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return cl.FeeForTx(ctx, tx)
}

func (cl *nodeClient) FeeForTx(ctx context.Context, tx *proto.TransferV2) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Debug("call service method 'FeeForTx'")
	url := cl.nodeClient.GetOptions().BaseUrl + calculateFeeUrl
	txJson, err := json.Marshal(tx)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(txJson))
	if err != nil {
		return 0, err
	}
	fee := new(FeeResponse)
	if _, err := cl.nodeClient.Do(ctx, req, fee); err != nil {
		log.Errorf("get fee fails: %s", err)
		return 0, err
	}
	return fee.FeeAmount, nil
}
