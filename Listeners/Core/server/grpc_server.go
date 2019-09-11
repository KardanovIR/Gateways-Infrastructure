package server

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
	"time"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/repositories"
)

type IGrpcServer interface {
	AddTask(ctx context.Context, in *pb.AddTaskRequest) (*pb.AddTaskResponse, error)
	RemoveTask(ctx context.Context, in *pb.RemoveTaskRequest) (*pb.Empty, error)
	RemoveTaskByTxHash(ctx context.Context, in *pb.RemoveTaskByTxHashRequest) (*pb.Empty, error)
}

type grpcServer struct {
	port      string
	rp        repositories.IRepository
	chainType models.ChainType
}

var (
	server                  IGrpcServer
	onceGrpcServertInstance sync.Once
)

func (s *grpcServer) AddTask(ctx context.Context, in *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	log := logger.FromContext(ctx)
	log.Infof("AddTask %+v", in)

	taskType, err := strconv.Atoi(in.TaskType)
	if err != nil {
		log.Errorf("Wrong task type %s: %s", in.TaskType, err)
		return nil, err
	}
	var newTask = models.Task{
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ListenTo:       models.ListenObject{Type: models.ListenType(in.ListenTo.Type), Value: in.ListenTo.Value},
		Callback:       models.Callback{Type: models.CallbackType(in.CallbackType), ProcessId: in.ProcessId},
		BlockchainType: s.chainType,
		Type:           models.TaskType(taskType),
	}

	id, err := s.rp.PutTask(ctx, &newTask)
	if err != nil {
		log.Errorf("Creating task fails: %s", err)
		return nil, err
	}

	return &pb.AddTaskResponse{TaskId: id}, nil
}

func (s *grpcServer) RemoveTask(ctx context.Context, in *pb.RemoveTaskRequest) (*pb.Empty, error) {
	log := logger.FromContext(ctx)
	log.Info("RemoveTask")

	var err = s.rp.RemoveTask(ctx, in.TaskId)
	if err != nil {
		log.Errorf("Removing task fails: %s", err)
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *grpcServer) RemoveTaskByTxHash(ctx context.Context, in *pb.RemoveTaskByTxHashRequest) (*pb.Empty, error) {
	log := logger.FromContext(ctx)
	log.Infof("RemoveTaskByTxHash %s", in.Hash)

	tasks, err := s.rp.FindByAddressOrTxId(ctx, s.chainType, "", in.Hash)
	if err != nil {
		log.Errorf("Removing task fails: %s", err)
		return nil, err
	}

	for _, task := range tasks {
		err = s.rp.RemoveTask(ctx, task.Id.Hex())
		if err != nil {
			log.Errorf("Removing task fails: %s", err)
			return nil, err
		}
	}

	return &pb.Empty{}, nil
}

func InitAndStart(ctx context.Context, port string, db repositories.IRepository, chainType models.ChainType) error {
	log := logger.FromContext(ctx)
	var initErr error
	onceGrpcServertInstance.Do(func() {
		server = &grpcServer{rp: db, port: ":" + port, chainType: chainType}

		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Errorf("failed to listen: %v", err)
			initErr = err
			return
		}

		newServer := grpc.NewServer()
		pb.RegisterListenerServer(newServer, server)
		log.Info("Grpc server registered")
		if err := newServer.Serve(lis); err != nil {
			log.Errorf("failed to serve: %v", err)
			initErr = err
			return
		}
	})

	return initErr
}

func GetGrpsServer() IGrpcServer {
	onceGrpcServertInstance.Do(func() {
		panic("try to get grpc server before it's creation!")
	})
	return server
}
