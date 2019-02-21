package server

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"

	pb "github.com/Waves/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/Waves/GatewaysInfrastructure/Adapters/Eth/services"
)

type IGrpcServer interface {
	GasPrice(ctx context.Context, in *pb.GasPriceRequest) (*pb.GasPriceReply, error)
}

type grpcServer struct {
	port       string
	nodeClient services.INodeClient
}

var (
	server                  IGrpcServer
	onceGrpcServertInstance sync.Once
)

func (s *grpcServer) GasPrice(ctx context.Context, in *pb.GasPriceRequest) (*pb.GasPriceReply, error) {
	log.Printf("Received")

	var gasPrice, err = s.nodeClient.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("connect to etherium node fails: ", err)
		return nil, err
	}

	return &pb.GasPriceReply{GasPrice: gasPrice.String()}, nil
}

func InitAndStart(port string, client services.INodeClient) error {
	var initErr error
	onceGrpcServertInstance.Do(func() {
		server = &grpcServer{nodeClient: client, port: ":" + port}

		lis, err := net.Listen("tcp", ":" + port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			initErr = err
			return
		}

		newServer := grpc.NewServer()
		pb.RegisterCommonServer(newServer, server)
		log.Printf("Grpc server registered")
		if err := newServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
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
