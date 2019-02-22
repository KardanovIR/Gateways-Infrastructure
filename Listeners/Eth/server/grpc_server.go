package server

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"sync"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/services"
)

type IGrpcServer interface {
	GasPrice(ctx context.Context, in *pb.GasPriceRequest) (*pb.GasPriceReply, error)
}

type grpcServer struct {
	port       string
	nodeClient services.INodeReader
}

var (
	server                  IGrpcServer
	onceGrpcServertInstance sync.Once
)

func (s *grpcServer) GasPrice(ctx context.Context, in *pb.GasPriceRequest) (*pb.GasPriceReply, error) {
	log := logger.FromContext(ctx)
	log.Info("GasPrice")

	var gasPrice, err = s.nodeClient.SuggestGasPrice(ctx)
	if err != nil {
		log.Errorf("connect to etherium node fails: %s", err)
		return nil, err
	}

	return &pb.GasPriceReply{GasPrice: gasPrice.String()}, nil
}

func InitAndStart(ctx context.Context, port string, client services.INodeClient) error {
	log := logger.FromContext(ctx)
	var initErr error
	onceGrpcServertInstance.Do(func() {
		server = &grpcServer{nodeClient: client, port: ":" + port}

		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Errorf("failed to listen: %v", err)
			initErr = err
			return
		}

		newServer := grpc.NewServer()
		pb.RegisterCommonServer(newServer, server)
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
