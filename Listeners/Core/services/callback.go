package services

import (
	"context"
	"google.golang.org/grpc"
	"sync"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/grpc/client"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
)

var (
	callbackServiceOnce sync.Once
	service             ICallbackService
)

type ICallbackService interface {
	SendRequest(ctx context.Context, task *models.Task, txId string) error
}

func NewCallbackService(ctx context.Context, callbackUrl string) error {
	log := logger.FromContext(ctx)
	var err error
	callbackServiceOnce.Do(func() {
		conn, e := grpc.Dial(callbackUrl, grpc.WithInsecure())
		if e != nil {
			err = e
			return
		}
		service = callbackService{pb.NewCoreServiceClient(conn)}
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
}

func (cs callbackService) SendRequest(ctx context.Context, task *models.Task, txId string) error {
	log := logger.FromContext(ctx)
	log.Debugf("send callback request %s for processId %s (txId %s)", task.Callback.Type, task.Callback.ProcessId, txId)
	requestFunc := cs.getCallbackByType(ctx, task.Callback.Type)
	_, err := requestFunc(ctx, &pb.Request{TxHash: txId, ProcessId: task.Callback.ProcessId})
	return err
}

func (cs callbackService) getCallbackByType(ctx context.Context, callbackType models.CallbackType) func(ctx context.Context, in *pb.Request, opts ...grpc.CallOption) (*pb.Empty, error) {
	switch callbackType {
	case models.InitOutTx:
		return cs.grpcClient.InitOutTx
	case models.FinishProcess:
		return cs.grpcClient.FinishProcess
	default:
		return func(ctx context.Context, in *pb.Request, opts ...grpc.CallOption) (*pb.Empty, error) {
			log := logger.FromContext(ctx)
			log.Errorf("not implemented callback %s was requested", callbackType)
			return &pb.Empty{}, nil
		}
	}
}
