package server

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"sync"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Router/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Router/service"
)

type grpcServer struct {
	port    string
	service service.IBlockChainsService
}

var (
	onceGrpcServerInstance sync.Once
)

func InitAndStart(ctx context.Context, port string, bs service.IBlockChainsService) error {
	log := logger.FromContext(ctx)
	var initErr error
	onceGrpcServerInstance.Do(func() {
		server := &grpcServer{service: bs, port: ":" + port}

		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Errorf("failed to listen: %v", err)
			initErr = err
			return
		}

		newServer := grpc.NewServer()
		pb.RegisterRouterServer(newServer, server)
		log.Infof("Grpc server registered on %s", port)
		if err := newServer.Serve(lis); err != nil {
			log.Errorf("failed to serve: %v", err)
			initErr = err
			return
		}
	})

	return initErr
}
