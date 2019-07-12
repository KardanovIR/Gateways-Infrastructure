package listener_test

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/grpc/client"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/repositories"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/server"
	coreServices "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/services"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Ergo/config"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Ergo/services"
	"google.golang.org/grpc"
)

const (
	txHashForSearchingByTxHash  = "1dd0dec03c1d3a6a96977a2f69d69e837975530a5825ee62485aa80f52527fca"
	txHashForSearchingByAddress = "4e735fc1f2beb14c479dd87b89c772b531f1ad699bdac07236a98c13706b531a"
	addressForSearching         = "3WzQunyAcbZVqCjzVDvDNHMQoSLopXzknYpKjsqz3LnWJgyFhHcY"
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

	mongoDB := mongoConnect(ctx, config.Cfg.Db.Host, config.Cfg.Db.Name)
	defer func() {
		if _, err := mongoDB.Collection(repositories.CChainState).DeleteOne(ctx, bson.D{{"chaintype", models.Ergo}}); err != nil {
			log.Error(err)
		}
	}()

	// add tasks Transfer
	_, err = grpcClient.AddTask(ctx,
		&blockchain.AddTaskRequest{
			ListenTo:     &blockchain.ListenObject{Type: "Address", Value: addressForSearching},
			CallbackType: string(models.InitInTx),
			ProcessId:    "111111111",
			TaskType:     strconv.Itoa(int(models.OneTime))},
	)

	if err != nil {
		log.Error("adding task fails", err)
		t.FailNow()
	}
	// add tasks Transfer
	_, err = grpcClient.AddTask(ctx,
		&blockchain.AddTaskRequest{
			ListenTo:     &blockchain.ListenObject{Type: "TxId", Value: txHashForSearchingByTxHash},
			CallbackType: string(models.FinishProcess),
			ProcessId:    "22222",
			TaskType:     strconv.Itoa(int(models.OneTime))},
	)

	if err != nil {
		log.Error("adding task fails", err)
		t.FailNow()
	}
	log.Debug("adding task success")

	err = services.GetNodeReader().Start(ctx)
	if err != nil {
		log.Error("node reader start fails", err)
		t.Fail()
		return
	}

	log.Debugf("node reader start  success", err)

	// wait for receiving callback
	var InitInTx bool
	var finishProcess bool

	for i := 0; i < 2; i++ {
		select {
		case callback := <-callBackChannel:
			if callback == "InitInTx" {
				InitInTx = true
			}
			if callback == "FinishProcess" {
				finishProcess = true
			}
		case <-time.After(10 * time.Second):
			log.Error("so long waiting...")
			t.FailNow()
		}
	}
	assert.True(t, InitInTx)
	assert.True(t, finishProcess)

	services.GetNodeReader().Stop(ctx)
}

func beforeTests(ctx context.Context, t *testing.T) {
	log, _ := logger.Init(false, logger.DEBUG)
	initTestOnce.Do(func() {
		err := config.Load("./testdata/config_test.yml")
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("Cfg: %+v", config.Cfg)
		if err := repositories.New(ctx, config.Cfg.Db.Host, config.Cfg.Db.Name); err != nil {
			log.Fatal("Can't create db connection: ", err)
		}

		if err := coreServices.NewCallbackService(ctx, config.Cfg.CallbackUrl, config.Cfg.Node.ChainType); err != nil {
			log.Fatal("Can't create callback service: ", err)
		}

		nodeClient := services.NewNodeClient(ctx, config.Cfg.Node)
		repository := repositories.GetRepository()
		if err := services.NewReader(ctx, &config.Cfg.Node, repository, nodeClient, coreServices.GetCallbackService()); err != nil {
			log.Fatal("Can't create node's client: ", err)
		}

		go func() {
			// start mock of core server to handle callbacks
			upCoreServer(t, ctx, config.Cfg.CallbackUrl)
		}()

		go func() {
			// start grpc server
			if err := server.InitAndStart(ctx, config.Cfg.Port, repositories.GetRepository(), config.Cfg.Node.ChainType); err != nil {
				log.Fatal("Can't start grpc server", err)
			}
		}()
		time.Sleep(100 * time.Millisecond)
	})

}

func getGrpcClient() (blockchain.ListenerClient, error) {
	conn, err := grpc.Dial(fmt.Sprint(":", config.Cfg.Port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return blockchain.NewListenerClient(conn), nil
}

// callbacks will be sent to this server
func upCoreServer(t *testing.T, ctx context.Context, url string) {
	log := logger.FromContext(ctx)
	lis, err := net.Listen("tcp", url)
	if err != nil {
		log.Fatal("failed to listen: %s", err)
	}
	newServer := grpc.NewServer()

	core.RegisterCoreServiceServer(newServer, coreServerMock{requestChannel: callBackChannel, t: t})
	if err := newServer.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}

type coreServerMock struct {
	requestChannel chan string
	t              *testing.T
}

func (c coreServerMock) InitInTx(ctx context.Context, in *core.InitInTxRequest) (*core.Empty, error) {
	assert.Equal(c.t, txHashForSearchingByAddress, in.TxHash)
	assert.Equal(c.t, addressForSearching, in.Address)
	callBackChannel <- "InitInTx"
	return &core.Empty{}, nil
}

func (c coreServerMock) CompleteTx(context.Context, *core.TxRequest) (*core.Empty, error) {
	panic("implement me")
}

func (c coreServerMock) InitOutTx(ctx context.Context, in *core.Request) (*core.Empty, error) {
	return nil, nil
}

func (c coreServerMock) FinishProcess(ctx context.Context, in *core.Request) (*core.Empty, error) {
	assert.Equal(c.t, "22222", in.ProcessId)
	assert.Equal(c.t, txHashForSearchingByTxHash, in.TxHash)
	callBackChannel <- "FinishProcess"
	return &core.Empty{}, nil
}

func mongoConnect(ctx context.Context, url string, dbName string) *mongo.Database {
	log := logger.FromContext(ctx)
	mongoClient, err := mongo.Connect(ctx, "mongodb://"+url)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB at %s: %v", url, err)
	}
	return mongoClient.Database(dbName)
}
