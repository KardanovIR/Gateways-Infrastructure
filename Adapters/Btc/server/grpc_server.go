package server

import (
	"context"
	dataClient "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/services/data_client"
	nodeClient "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/services/node_client"
	"google.golang.org/grpc"
	"net"
	"sync"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc/logger"
)

type grpcServer struct {
	port       string
	nodeClient nodeClient.INodeClient
	dataClient dataClient.IDataClient
}

var (
	server                  pb.AdapterServer
	onceGrpcServertInstance sync.Once
)

func InitAndStart(ctx context.Context, port string, client nodeClient.INodeClient, dataClient dataClient.IDataClient) error {
	log := logger.FromContext(ctx)
	var initErr error
	onceGrpcServertInstance.Do(func() {
		server = &grpcServer{nodeClient: client, port: ":" + port, dataClient: dataClient}

		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Errorf("failed to listen: %v", err)
			initErr = err
			return
		}

		newServer := grpc.NewServer()
		pb.RegisterAdapterServer(newServer, server)
		log.Info("Grpc server registered")
		if err := newServer.Serve(lis); err != nil {
			log.Errorf("failed to serve: %v", err)
			initErr = err
			return
		}
	})

	return initErr
}

func GetGrpsServer() pb.AdapterServer {
	onceGrpcServertInstance.Do(func() {
		panic("try to get grpc server before it's creation!")
	})
	return server
}
