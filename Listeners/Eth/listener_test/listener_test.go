package listener_test

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	pb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/grpc"
	corePb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/grpc/client"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/repositories"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/server"
	coreServices "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/services"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CoreServerPort = "5020"
)

var (
	callBackChannel = make(chan string, 2)
	initTestOnce    sync.Once
)

// need Parity node to read internal tx: add NODE_PARITY_HOST parameter to env variable
func TestListenerEth(t *testing.T) {
	ctx := context.Background()
	// setup
	beforeTests(ctx, t, 5202062)
	log := logger.FromContext(ctx)
	grpcClient, err := getGrpcClient()
	if err != nil {
		log.Fatal("Can't init grpc client", err)
	}

	// add tasks Transfer
	_, err = grpcClient.AddTask(ctx,
		&pb.AddTaskRequest{
			ListenTo:     &pb.ListenObject{Type: "Address", Value: "0x9515735d60E8fF4036EFAFFAf3370F3097615d19"},
			CallbackType: string(models.InitOutTx),
			ProcessId:    "111111111",
			TaskType:     strconv.Itoa(int(models.OneTime))},
	)
	if err != nil {
		log.Error("adding task fails", err)
		t.FailNow()
		return
	}
	// internal in block 5202068
	_, err = grpcClient.AddTask(ctx,
		&pb.AddTaskRequest{
			ListenTo:     &pb.ListenObject{Type: "Address", Value: "0x50F554649ED757D40d5Bd32B1154AFfc4278359B"},
			CallbackType: string(models.InitInTx),
			ProcessId:    "",
			TaskType:     strconv.Itoa(int(models.OneTime))},
	)

	if err != nil {
		log.Error("adding task fails", err)
		t.FailNow()
		return
	}
	log.Debugf("adding task success", err)

	err = services.GetNodeReader().Start(ctx)
	if err != nil {
		log.Error("node reader start fails", err)
		t.Fail()
		return
	}

	log.Debugf("node reader start  success", err)

	// wait for receiving callback
	var isTransfer bool
	var isInternalTransfer bool

	for i := 0; i < 2; i++ {
		select {
		case callback := <-callBackChannel:
			if callback == "InitOutTx" {
				isTransfer = true
			}
			if callback == "InitInTx" {
				isInternalTransfer = true
			}
		case <-time.After(30 * time.Second):
			log.Error("so long waiting...")
			t.FailNow()
		}
	}
	assert.True(t, isTransfer)
	assert.True(t, isInternalTransfer)

	mongoDB := mongoConnect(ctx, config.Cfg.Db.Host, config.Cfg.Db.Name)
	defer func() {
		if _, err := mongoDB.Collection(repositories.CChainState).DeleteOne(ctx, bson.D{{"chaintype", models.Ethereum}}); err != nil {
			log.Error(err)
		}
	}()

	services.GetNodeReader().Stop(ctx)
}

func TestListenerErc20(t *testing.T) {
	ctx := context.Background()
	// setup
	beforeTests(ctx, t, 6197315)
	log := logger.FromContext(ctx)
	grpcClient, err := getGrpcClient()
	if err != nil {
		log.Fatal("Can't init grpc client", err)
	}
	mongoDB := mongoConnect(ctx, config.Cfg.Db.Host, config.Cfg.Db.Name)
	defer func() {
		if _, err := mongoDB.Collection(repositories.CChainState).DeleteOne(ctx, bson.D{{"chaintype", models.Ethereum}}); err != nil {
			log.Error(err)
		}
	}()
	// add tasks Transfer
	_, err = grpcClient.AddTask(ctx,
		&pb.AddTaskRequest{
			ListenTo:     &pb.ListenObject{Type: "Address", Value: "0x8ec23aCbe3Eed99E92d6D7a85a27A45dA3A04e7d"},
			CallbackType: string(models.InitInTx),
			TaskType:     strconv.Itoa(int(models.OneTime))},
	)

	if err != nil {
		log.Error("adding task fails", err)
		t.FailNow()
		return
	}
	log.Debugf("adding task success", err)

	// in 6197316 has failed tx for 0x8ec23aCbe3Eed99E92d6D7a85a27A45dA3A04e7d
	// in 6197336 has success tx for 0x8ec23aCbe3Eed99E92d6D7a85a27A45dA3A04e7d
	err = services.GetNodeReader().Start(ctx)
	if err != nil {
		log.Error("node reader start fails", err)
		t.Fail()
		return
	}

	log.Debugf("node reader start  success", err)

	// wait for receiving callback
	var isTransfer bool

	select {
	case callback := <-callBackChannel:
		log.Info("!!!!!!!!!!!!!", callback)
		if callback == "InitInTx" {
			isTransfer = true
		}
	case <-time.After(300 * time.Second):
		log.Error("so long waiting...")
		t.FailNow()
	}

	if !isTransfer {
		t.Fail()
	}

	services.GetNodeReader().Stop(ctx)
}

func beforeTests(ctx context.Context, t *testing.T, startBlock int64) {
	log, _ := logger.Init(false, logger.DEBUG)
	initTestOnce.Do(func() {
		err := config.Load("./testdata/config_test.yml")
		if err != nil {
			log.Fatal(err)
		}
		config.Cfg.Node.StartBlockHeight = startBlock

		if err := repositories.New(ctx, config.Cfg.Db.Host, config.Cfg.Db.Name); err != nil {
			log.Fatal("Can't create db connection: ", err)
		}

		// create node reader
		if err = services.New(ctx, &config.Cfg.Node, repositories.GetRepository()); err != nil {
			log.Fatal(err)
		}

		if err := coreServices.NewCallbackService(ctx, config.Cfg.CallbackUrl, models.Ethereum); err != nil {
			log.Fatal("Can't create callback service: ", err)
		}

		go func() {
			// start mock of core server to handle callbacks
			upCoreServer(t, ctx, CoreServerPort)
		}()

		go func() {
			// start grpc server
			if err := server.InitAndStart(ctx, config.Cfg.Port, repositories.GetRepository(), models.Ethereum); err != nil {
				log.Fatal("Can't start grpc server", err)
			}
		}()
		time.Sleep(100 * time.Millisecond)
	})

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

func (c coreServerMock) InitInTx(ctx context.Context, in *corePb.InitInTxRequest) (*corePb.Empty, error) {
	if in.Address == "0x8ec23aCbe3Eed99E92d6D7a85a27A45dA3A04e7d" {
		// erc-20 test
		assert.Equal(c.t, "0x8ec23aCbe3Eed99E92d6D7a85a27A45dA3A04e7d", in.Address)
		assert.Equal(c.t, "0x9552c6303ae43bd9b4d96bd31eca00faac6abe9c68511b8591ca74c588bb1e52", in.TxHash)
	} else {
		// internal
		assert.Equal(c.t, "0x50F554649ED757D40d5Bd32B1154AFfc4278359B", in.Address)
		assert.Equal(c.t, "0xa874bd3d557d5145f71a354c8d4035acc05d7778b2c26c2e73d2621cd8bab143", in.TxHash)
	}
	callBackChannel <- "InitInTx"
	return &corePb.Empty{}, nil
}

func (c coreServerMock) CompleteTx(context.Context, *corePb.TxRequest) (*corePb.Empty, error) {
	panic("implement me")
}

func (c coreServerMock) InitOutTx(ctx context.Context, in *corePb.Request) (*corePb.Empty, error) {
	assert.Equal(c.t, "111111111", in.ProcessId)
	assert.Equal(c.t, "0x48094cb5687722d386909162257fa3b07c74e9fca5d03c38331fbcd818544f5b", in.TxHash)
	callBackChannel <- "InitOutTx"
	return &corePb.Empty{}, nil
}

func (coreServerMock) FinishProcess(context.Context, *corePb.Request) (*corePb.Empty, error) {
	callBackChannel <- "FinishProcess"
	return &corePb.Empty{}, nil
}

func mongoConnect(ctx context.Context, url string, dbName string) *mongo.Database {
	log := logger.FromContext(ctx)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+url))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB at %s: %v", url, err)
	}
	return mongoClient.Database(dbName)
}
