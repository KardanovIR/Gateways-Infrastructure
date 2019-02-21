package server

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/services"
)

type IGrpcServer interface {
	GasPrice(ctx context.Context, in *pb.GasPriceRequest) (*pb.GasPriceReply, error)
	Start() error
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

func (s *grpcServer) Start() error {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	newServer := grpc.NewServer()
	pb.RegisterCommonServer(newServer, &grpcServer{})
	log.Printf("Grpc server registered")
	if err := newServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}

func New(port string, client services.INodeClient) error {
	var err error
	onceGrpcServertInstance.Do(func() {
		server = &grpcServer{nodeClient: client, port: ":" + port}
	})
	return err
}

func GetGrpsServer() IGrpcServer {
	onceGrpcServertInstance.Do(func() {
		panic("try to get grpc server before it's creation!")
	})
	return server
}
