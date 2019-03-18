package server

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"

	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/config"
	pb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/repositories"
)

type IGrpcServer interface {
	AddTask(ctx context.Context, in *pb.AddTaskRequest) (*pb.AddTaskResponse, error)
	RemoveTask(ctx context.Context, in *pb.RemoveTaskRequest) (*pb.Empty, error)
}

type grpcServer struct {
	port string
	rp   repositories.IRepository
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
		Address:        in.Address,
		Callback:       models.Callback{in.CallbackUrl, models.CallbackType(in.CallbackType), nil},
		BlockchainType: config.Cfg.Node.ChainType,
		Type:           models.TaskType(taskType),
	}

	id, err := s.rp.PutTask(ctx, newTask)
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

func InitAndStart(ctx context.Context, port string, db repositories.IRepository) error {
	log := logger.FromContext(ctx)
	var initErr error
	onceGrpcServertInstance.Do(func() {
		server = &grpcServer{rp: db, port: ":" + port}

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