package server

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/config"
	pb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
	"google.golang.org/grpc"
)

const serverPort = "20001"
const addedTaskId = "123"

var oneServer sync.Once

func TestGrpcServerAddTask(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	conn, err := startServerAndGetConnection(ctx)
	if err != nil {
		log.Error("connection to grpc server don't establish", err)
		t.Fail()
	}
	defer conn.Close()

	c := pb.NewListenerClient(conn)
	reply, err := c.AddTask(ctx, &pb.AddTaskRequest{
		ListenTo: &pb.ListenObject{Type: "Address", Value: "123456"}, CallbackType: string(models.Post), TaskType: "1"})
	if err != nil {
		log.Error("adding task fails", err)
		t.Fail()
	}
	if reply.TaskId != addedTaskId {
		t.Fail()
	}
}

func TestGrpcServerRemoveTask(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	conn, err := startServerAndGetConnection(ctx)
	if err != nil {
		log.Error("connection to grpc server don't establish", err)
		t.Fail()
	}
	defer conn.Close()

	c := pb.NewListenerClient(conn)
	_, err = c.RemoveTask(ctx, &pb.RemoveTaskRequest{TaskId: addedTaskId})
	if err != nil {
		log.Error("task's removing fails", err)
		t.Fail()
	}
}

func startServerAndGetConnection(ctx context.Context) (*grpc.ClientConn, error) {
	log := logger.FromContext(ctx)
	config.Cfg = &config.Config{Node: config.Node{ChainType: models.Ethereum}}
	oneServer.Do(func() {
		go func() {
			err := InitAndStart(ctx, serverPort, &repoMock{})
			if err != nil {
				log.Fatal("server can't start", err)
			}
		}()
	})
	return grpc.Dial(fmt.Sprint(":", serverPort), grpc.WithInsecure())
}

type repoMock struct {
}

func (n *repoMock) PutTask(ctx context.Context, task *models.Task) (string, error) {
	return addedTaskId, nil
}

func (n *repoMock) RemoveTask(ctx context.Context, id string) error {
	return nil
}

func (n *repoMock) FindByAddress(ctx context.Context, ticket models.ChainType, addresses string) (tasks []*models.Task, err error) {
	return make([]*models.Task, 0), nil
}

func (n *repoMock) FindByAddressOrTxId(ctx context.Context, ticket models.ChainType, address string, txID string) (tasks []*models.Task, err error) {
	return make([]*models.Task, 0), nil
}

func (n *repoMock) GetLastChainState(ctx context.Context, chainType models.ChainType) (chainState *models.ChainState, err error) {
	return new(models.ChainState), nil
}

func (n *repoMock) PutChainState(ctx context.Context, state *models.ChainState) (newState *models.ChainState, err error) {
	return &models.ChainState{}, nil
}
