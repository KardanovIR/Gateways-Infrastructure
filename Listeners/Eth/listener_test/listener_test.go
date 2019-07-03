package listener_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/config"
	pb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/grpc"
	corePb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/grpc/client"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/repositories"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/services"
)

const (
	CoreServerPort = "5020"
)

var (
	callBackChannel = make(chan string, 2)
	initTestOnce    sync.Once
)

func TestListener(t *testing.T) {
	ctx := context.Background()
	// setup
	beforeTests(ctx, t)
	log := logger.FromContext(ctx)
	grpcClient, err := getGrpcClient()
	if err != nil {
		log.Fatal("Can't init grpc client", err)
	}

	// add tasks Transfer
	_, err = grpcClient.AddTask(ctx,
		&pb.AddTaskRequest{
			ListenTo:     &pb.ListenObject{Type: "Address", Value: "0x9515735d60e8ff4036efaffaf3370f3097615d19"},
			CallbackType: string(models.InitOutTx),
			ProcessId:    "111111111",
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

	select {
	case callback := <-callBackChannel:
		if callback == "InitOutTx" {
			isTransfer = true
		}
	case <-time.After(10 * time.Second):
		log.Error("so long waiting...")
		t.FailNow()
	}

	if !isTransfer {
		t.Fail()
	}

	mongoDB := mongoConnect(ctx, config.Cfg.Db.Host, config.Cfg.Db.Name)
	defer func() {
		if _, err := mongoDB.Collection(repositories.CChainState).DeleteOne(ctx, bson.D{{"chaintype", models.Ethereum}}); err != nil {
			log.Error(err)
		}
	}()

	services.GetNodeReader().Stop(ctx)
}

func beforeTests(ctx context.Context, t *testing.T) {
	log, _ := logger.Init(false, logger.DEBUG)
	initTestOnce.Do(func() {
		err := config.Load("./testdata/config_test.yml")
		if err != nil {
			log.Fatal(err)
		}

		if err := repositories.New(ctx, config.Cfg.Db.Host, config.Cfg.Db.Name); err != nil {
			log.Fatal("Can't create db connection: ", err)
		}

		// create node reader
		if err = services.New(ctx, &config.Cfg.Node, repositories.GetRepository()); err != nil {
			log.Fatal(err)
		}

		if err := services.NewCallbackService(ctx, config.Cfg.CallbackUrl); err != nil {
			log.Fatal("Can't create callback service: ", err)
		}

		go func() {
			// start mock of core server to handle callbacks
			upCoreServer(t, ctx, CoreServerPort)
		}()

		go func() {
			// start grpc server
			if err := server.InitAndStart(ctx, config.Cfg.Port, repositories.GetRepository()); err != nil {
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
	mongoClient, err := mongo.Connect(ctx, "mongodb://"+url)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB at %s: %v", url, err)
	}
	return mongoClient.Database(dbName)
}
