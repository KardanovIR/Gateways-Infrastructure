package server

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"sync"

	_ "github.com/jnewmano/grpc-json-proxy/codec"
	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/server/converter"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/services"
)

type grpcServer struct {
	port       string
	nodeClient services.INodeClient
	converter  converter.IConverter
}

var (
	server                  pb.AdapterServer
	onceGrpcServertInstance sync.Once
)

func InitAndStart(ctx context.Context, port string, client services.INodeClient, converter converter.IConverter) error {
	log := logger.FromContext(ctx)
	var initErr error
	onceGrpcServertInstance.Do(func() {
		server = &grpcServer{nodeClient: client, port: ":" + port, converter: converter}

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
