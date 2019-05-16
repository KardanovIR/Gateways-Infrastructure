package listener_test

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/config"
	pb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/grpc"
	corePb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/grpc/client"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/repositories"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/services"
)

const (
	httpServerPort = "5020"
)

var callBackChannel = make(chan string, 5)

func TestListener(t *testing.T) {
	ctx := context.Background()
	// setup
	beforeTests(t, ctx)
	log := logger.FromContext(ctx)
	grpcClient, err := getGrpcClient()
	if err != nil {
		log.Fatal("Can't init grpc client", err)
	}
	mongoClient := mongoConnect(ctx, config.Cfg.Db.Host, config.Cfg.Db.Name)

	// add tasks TransferV1 address
	_, err = grpcClient.AddTask(ctx,
		&pb.AddTaskRequest{
			ListenTo:     &pb.ListenObject{Type: "Address", Value: "3PAgyfDELn1UixKCLQ6UsVakuofXXZMdYC4"},
			CallbackType: string(models.InitOutTx),
			TaskType:     strconv.Itoa(int(models.OneTime)),
			ProcessId:    "111111111",
		})

	if err != nil {
		log.Error("adding task fails", err)
		t.FailNow()
		return
	}
	// add tasks TransferV1 TxID
	_, err = grpcClient.AddTask(ctx,
		&pb.AddTaskRequest{
			ListenTo:     &pb.ListenObject{Type: "TxId", Value: "1tASvqX4TVNARYZZ7w1JnwUBr8pXtXHPiRVYMARYVRJ"},
			CallbackType: string(models.InitOutTx),
			TaskType:     strconv.Itoa(int(models.OneTime)),
			ProcessId:    "222222222",
		})

	if err != nil {
		log.Error("adding task fails", err)
		t.FailNow()
		return
	}

	// add tasks TransferV2
	_, err = grpcClient.AddTask(ctx,
		&pb.AddTaskRequest{
			ListenTo:     &pb.ListenObject{Type: "Address", Value: "3P63utQnWvQ7Xd8NMVFYjd1UqrUBqXsFVr8"},
			CallbackType: string(models.FinishProcess),
			TaskType:     strconv.Itoa(int(models.OneTime)),
			ProcessId:    "333333333"},
	)

	if err != nil {
		log.Error("adding task fails", err)
		t.FailNow()
		return
	}

	// add tasks Mass Transfer
	_, err = grpcClient.AddTask(ctx,
		&pb.AddTaskRequest{
			ListenTo:     &pb.ListenObject{Type: "Address", Value: "3PAc93kp7CDwh2tc682JqDKT96uP5XeHsH2"},
			CallbackType: string(models.FinishProcess),
			TaskType:     strconv.Itoa(int(models.OneTime)),
			ProcessId:    "444444444",
		})

	if err != nil {
		log.Error("adding task fails", err)
		t.FailNow()
		return
	}

	err = services.GetNodeReader().Start(ctx)
	if err != nil {
		log.Error("node reader start fails", err)
		t.Fail()
		return
	}
	defer func() {
		_, err := mongoClient.Collection(repositories.CChainState).DeleteOne(ctx, bson.D{{"chaintype", models.Waves}})
		if err != nil {
			log.Error("node reader: clearing test fails", err)
		}
	}()
	// wait for receiving callback
	var isTransfer, isTransfer2, isMassTransfer, isTaskByTxId bool
	for i := 0; i < 4; i++ {
		select {
		case callback := <-callBackChannel:
			if callback == "transfer" {
				isTransfer = true
			}
			if callback == "transfer2" {
				isTransfer2 = true
			}
			if callback == "masstransfer" {
				isMassTransfer = true
			}
			if callback == "transfer_txId" {
				isTaskByTxId = true
			}

		case <-time.After(10 * time.Second):
			log.Error("so long waiting...")
			t.FailNow()
		}
	}
	if !isTransfer {
		t.Fail()
	}
	if !isTransfer2 {
		t.Fail()
	}
	if !isMassTransfer {
		t.Fail()
	}
	if !isTaskByTxId {
		t.Fail()
	}
	services.GetNodeReader().Stop(ctx)
}

func beforeTests(t *testing.T, ctx context.Context) {
	log, _ := logger.Init(false, logger.DEBUG)
	err := config.Load("./testdata/config_test.yml")
	if err != nil {
		log.Fatal(err)
	}

	if err := repositories.New(ctx, config.Cfg.Db.Host, config.Cfg.Db.Name); err != nil {
		log.Fatal("Can't create db connection: ", err)
	}
	if err := services.NewCallbackService(ctx, config.Cfg.CallbackUrl, config.Cfg.Node.ChainType); err != nil {
		log.Fatal("Can't create callback service: ", err)
	}
	// create node reader
	if err = services.New(ctx, &config.Cfg.Node, repositories.GetRepository()); err != nil {
		log.Fatal(err)
	}

	go func() {
		// start http server to handle callbacks
		upCoreServer(t, ctx, httpServerPort)
	}()

	go func() {
		// start grpc server
		if err := server.InitAndStart(ctx, config.Cfg.Port, repositories.GetRepository()); err != nil {
			log.Fatal("Can't start grpc server", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

func getGrpcClient() (pb.ListenerClient, error) {
	conn, err := grpc.Dial(fmt.Sprint(":", config.Cfg.Port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewListenerClient(conn), nil
}

// callbacks will be sent to this server
func upCoreServer(t *testing.T, ctx context.Context, port string) {
	log := logger.FromContext(ctx)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("failed to listen: %s", err)
	}
	newServer := grpc.NewServer()
	corePb.RegisterCoreServiceServer(newServer, coreServerMock{requestChannel: callBackChannel, t: t})
	if err := newServer.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}

type coreServerMock struct {
	requestChannel chan string
	t              *testing.T
}

func (c coreServerMock) InitInTx(context.Context, *corePb.InitInTxRequest) (*corePb.Empty, error) {
	panic("implement me")
}

func (c coreServerMock) InitOutTx(ctx context.Context, in *corePb.Request) (*corePb.Empty, error) {
	assert.True(c.t, len(in.ProcessId) > 0)
	switch in.ProcessId {
	case "111111111":
		assert.Equal(c.t, "GA7PECC8DFEPwRwyN75FG2mpx5ad3BYAPSUp66MeT6RP", in.TxHash)
		callBackChannel <- "transfer"
	case "222222222":
		assert.Equal(c.t, "1tASvqX4TVNARYZZ7w1JnwUBr8pXtXHPiRVYMARYVRJ", in.TxHash)
		callBackChannel <- "transfer_txId"
	default:
		c.t.Fail()
	}

	return &corePb.Empty{}, nil
}

func (c coreServerMock) FinishProcess(ctx context.Context, in *corePb.Request) (*corePb.Empty, error) {
	assert.True(c.t, len(in.ProcessId) > 0)
	switch in.ProcessId {
	case "333333333":
		assert.Equal(c.t, "6BabdgUzv96pRcD42YPrv6Wd8oHAPbHhdPoMGTa2ziE9", in.TxHash)
		callBackChannel <- "transfer2"
	case "444444444":
		assert.Equal(c.t, "6bwsRivXTBU396bYswL6vW8rPYcbafKvbeAvM2sGKqLP", in.TxHash)
		callBackChannel <- "masstransfer"
	default:
		c.t.Fail()
	}
	return &corePb.Empty{}, nil
}

func mongoConnect(ctx context.Context, url string, dbName string) *mongo.Database {
	log := logger.FromContext(ctx)
	mongoClient, err := mongo.Connect(ctx, "mongodb://"+url)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB at %s: %v", url, err)
	}
	return mongoClient.Database(dbName)
}
