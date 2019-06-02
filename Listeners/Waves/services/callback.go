package services

import (
	"context"
	"google.golang.org/grpc"
	"sync"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/grpc/client"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/models"
)

var (
	callbackServiceOnce sync.Once
	service             ICallbackService
)

type ICallbackService interface {
	SendRequest(ctx context.Context, task *models.Task, txId string) error
}

func NewCallbackService(ctx context.Context, callbackUrl string, chainType models.ChainType) error {
	log := logger.FromContext(ctx)
	var err error
	callbackServiceOnce.Do(func() {
		conn, e := grpc.Dial(callbackUrl, grpc.WithInsecure())
		if e != nil {
			err = e
			return
		}
		service = callbackService{pb.NewCoreServiceClient(conn), chainType}
	})
	if err != nil {
		log.Errorf("error during initialize callback service: %s", err)
	}
	return err
}

func GetCallbackService() ICallbackService {
	callbackServiceOnce.Do(func() {
		panic("try to get callback service before it's creation!")
	})
	return service
}

type callbackService struct {
	grpcClient pb.CoreServiceClient
	chainType  models.ChainType
}

func (cs callbackService) SendRequest(ctx context.Context, task *models.Task, txId string) error {
	log := logger.FromContext(ctx)
	log.Debugf("send callback request %s for processId %s (txId %s)", task.Callback.Type, task.Callback.ProcessId, txId)
	_, err := cs.callback(ctx, txId, task.Callback.ProcessId, task.Callback.Type)
	return err
}

func (cs callbackService) callback(ctx context.Context, txHash string, processId string, callbackType models.CallbackType) (*pb.Empty, error) {
	switch callbackType {
	case models.InitOutTx:
		return cs.grpcClient.InitOutTx(ctx, &pb.Request{TxHash: txHash, ProcessId: processId})
	case models.FinishProcess:
		return cs.grpcClient.FinishProcess(ctx, &pb.Request{TxHash: txHash, ProcessId: processId})
	case models.InitInTx:
		return cs.grpcClient.InitInTx(ctx, &pb.TxRequest{TxHash: txHash, BlockchainType: string(cs.chainType)})
	case models.CompleteTx:
		return cs.grpcClient.CompleteTx(ctx, &pb.TxRequest{TxHash: txHash, BlockchainType: string(cs.chainType)})
	default:
		log := logger.FromContext(ctx)
		log.Errorf("not implemented callback %s was requested", callbackType)
		return &pb.Empty{}, nil
	}
}
