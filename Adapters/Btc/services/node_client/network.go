package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/converter"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
)

const (
	feeEstimateMode    = "ECONOMICAL"
	confTargetInBlocks = 1
)

type FeeResponse struct {
	FeeRate float64  `json:"feerate"`
	Blocks  int      `json:"blocks"`
	Errors  []string `json:"errors"`
}

func (cl *nodeClient) FeeRateForKByte(ctx context.Context) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'FeeRateForKByte'")
	targetBlocksByte, _ := json.Marshal(confTargetInBlocks)
	mode, _ := json.Marshal(feeEstimateMode)
	rawResult, err := cl.nodeClient.RawRequest("estimatesmartfee", []json.RawMessage{targetBlocksByte, mode})
	if err != nil {
		log.Error(err)
		return 0, err
	}
	log.Debugf("fee response: %s", string(rawResult))
	feeResponse := new(FeeResponse)
	err = json.Unmarshal(rawResult, feeResponse)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	if len(feeResponse.Errors) > 0 {
		for _, e := range feeResponse.Errors {
			log.Error(e)
		}
		return 0, fmt.Errorf("%+v", feeResponse.Errors)
	}
	fee, err := converter.GetIntFromFloat(feeResponse.FeeRate)
	log.Infof("current fee rate %d", fee)
	return fee, nil
}

func (cl *nodeClient) Fee(ctx context.Context, feeRate uint64, txSize int) uint64 {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'Fee' for fee rate %d, tx size %d", feeRate, txSize)
	var fee uint64
	m := feeRate * uint64(txSize)
	fee = m / 1000
	if m%1000 > 0 {
		fee += 1 // in previous step float amount was be truncated
	}
	log.Infof("fee %d", fee)
	return fee
}
