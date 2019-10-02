package services

import (
	"context"

	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/repository"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
)

type unspentTxService struct {
	rep repository.IUTXORepository
}

func NewUnspentTxService(ctx context.Context, rep repository.IUTXORepository) *unspentTxService {
	return &unspentTxService{rep}
}

// not thread safe
func (s *unspentTxService) addTxInputs(ctx context.Context, TxHash string, amount uint64, address string, outputN uint32) error {
	log := logger.FromContext(ctx)
	log.Infof("add input for address %s: amount %d, tx %s, out %d", address, amount, TxHash, outputN)
	unspent, err := s.rep.GetUnspentTxByTxHashAndOutputNumber(ctx, TxHash, outputN)
	if err != nil {
		return err
	}
	if unspent != nil {
		log.Warnf("output %s N %d is already added to db as unspent input", TxHash, outputN)
		return nil
	}
	return s.rep.AddUnspentTx(ctx, models.UnspentTx{Address: address, TxHash: TxHash, Amount: amount, OutputN: outputN})
}

func (s *unspentTxService) deleteTxInputs(ctx context.Context, txHash string, outN uint32) error {
	log := logger.FromContext(ctx)
	log.Debugf("delete input: hash %s, N %d", txHash, outN)
	u, err := s.rep.DeleteUnspentTx(ctx, txHash, outN)
	if err != nil {
		return err
	}
	if u == nil {
		log.Warnf("output with hash %s, N %d is not found", txHash, outN)
		return nil
	}
	if !u.Locked {
		log.Warnf("delete output with hash %s, N %d which was not locked", txHash, outN)
	}
	return nil
}
