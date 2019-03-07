package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"math/big"
	"net"
	"sync"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/services"
)

type IGrpcServer interface {
	GasPrice(ctx context.Context, in *pb.GasPriceRequest) (*pb.GasPriceReply, error)
	GetRawTransaction(ctx context.Context, in *pb.RawTransactionRequest) (*pb.RawTransactionReply, error)
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
	log := logger.FromContext(ctx)
	log.Info("GasPrice")

	var gasPrice, err = s.nodeClient.SuggestGasPrice(ctx)
	if err != nil {
		log.Errorf("connect to etherium node fails: %s", err)
		return nil, err
	}

	return &pb.GasPriceReply{GasPrice: gasPrice.String()}, nil
}

func (s *grpcServer) GetRawTransaction(ctx context.Context, in *pb.RawTransactionRequest) (*pb.RawTransactionReply, error) {
	log := logger.FromContext(ctx)
	log.Infof("getRawTransaction %+v", in)
	amount, ok := new(big.Int).SetString(in.Amount, 10)
	if !ok {
		err := fmt.Errorf("wrong amount value: %s", in.Amount)
		log.Error(err)
		return nil, err
	}
	var tx, err = s.nodeClient.CreateRawTransaction(ctx, in.AddressFrom, in.AddressTo, amount)
	if err != nil {
		log.Errorf("transaction's creation fails: %s", err)
		return nil, err
	}
	return &pb.RawTransactionReply{Tx: string(tx)}, nil
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
