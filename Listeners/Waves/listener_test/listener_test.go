package listener_test

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/config"
	pb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/repositories"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/services"
)

const (
	httpServerPort = "8085"
)

var callBackChannel = make(chan string, 2)

func TestListener(t *testing.T) {
	ctx := context.Background()
	// setup
	beforeTests(ctx)
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
			CallbackType: string(models.Get),
			CallbackUrl:  fmt.Sprintf("http://localhost:%s/transfer", httpServerPort),
			TaskType:     strconv.Itoa(int(models.OneTime))})

	if err != nil {
		log.Error("adding task fails", err)
		t.FailNow()
		return
	}
	// add tasks TransferV1 TxID
	_, err = grpcClient.AddTask(ctx,
		&pb.AddTaskRequest{
			ListenTo:     &pb.ListenObject{Type: "TxId", Value: "1tASvqX4TVNARYZZ7w1JnwUBr8pXtXHPiRVYMARYVRJ"},
			CallbackType: string(models.Get),
			CallbackUrl:  fmt.Sprintf("http://localhost:%s/transfer/txId", httpServerPort),
			TaskType:     strconv.Itoa(int(models.OneTime))})

	if err != nil {
		log.Error("adding task fails", err)
		t.FailNow()
		return
	}

	// add tasks TransferV2
	_, err = grpcClient.AddTask(ctx,
		&pb.AddTaskRequest{
			ListenTo:     &pb.ListenObject{Type: "Address", Value: "3P63utQnWvQ7Xd8NMVFYjd1UqrUBqXsFVr8"},
			CallbackType: string(models.Get),
			CallbackUrl:  fmt.Sprintf("http://localhost:%s/transfer2", httpServerPort),
			TaskType:     strconv.Itoa(int(models.OneTime))})

	if err != nil {
		log.Error("adding task fails", err)
		t.FailNow()
		return
	}

	// add tasks Mass Transfer
	_, err = grpcClient.AddTask(ctx,
		&pb.AddTaskRequest{
			ListenTo:     &pb.ListenObject{Type: "Address", Value: "3PAc93kp7CDwh2tc682JqDKT96uP5XeHsH2"},
			CallbackType: string(models.Get),
			CallbackUrl:  fmt.Sprintf("http://localhost:%s/masstransfer", httpServerPort),
			TaskType:     strconv.Itoa(int(models.OneTime))})

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

func beforeTests(ctx context.Context) {
	log, _ := logger.Init(false, logger.DEBUG)
	err := config.Load("./testdata/config_test.yml")
	if err != nil {
		log.Fatal(err)
	}
	if err := services.NewRestClient(ctx); err != nil {
		log.Fatal("Can't create rest client: ", err)
	}

	if err := repositories.New(ctx, config.Cfg.Db.Host, config.Cfg.Db.Name); err != nil {
		log.Fatal("Can't create db connection: ", err)
	}

	// create node reader
	if err = services.New(ctx, &config.Cfg.Node, services.GetRestClient(), repositories.GetRepository()); err != nil {
		log.Fatal(err)
	}

	go func() {
		// start http server to handle callbacks
		upHttpServer(ctx, httpServerPort)
	}()

	go func() {
		// start grpc server
		if err := server.InitAndStart(ctx, config.Cfg.Port, repositories.GetRepository()); err != nil {
			log.Fatal("Can't start grpc server", err)
		}
	}()
}

func getGrpcClient() (pb.ListenerClient, error) {
	conn, err := grpc.Dial(fmt.Sprint(":", config.Cfg.Port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewListenerClient(conn), nil
}

func upHttpServer(ctx context.Context, port string) {
	log := logger.FromContext(ctx)
	router := gin.Default()
	router.GET("/transfer", func(c *gin.Context) {
		callBackChannel <- "transfer"
		c.JSON(http.StatusOK, "")
	})
	router.GET("/transfer/txId", func(c *gin.Context) {
		callBackChannel <- "transfer_txId"
		c.JSON(http.StatusOK, "")
	})
	router.GET("/transfer2", func(c *gin.Context) {
		callBackChannel <- "transfer2"
		c.JSON(http.StatusOK, "")
	})
	router.GET("/masstransfer", func(c *gin.Context) {
		callBackChannel <- "masstransfer"
		c.JSON(http.StatusOK, "")
	})
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func mongoConnect(ctx context.Context, url string, dbName string) *mongo.Database {
	log := logger.FromContext(ctx)
	mongoClient, err := mongo.Connect(ctx, "mongodb://"+url)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB at %s: %v", url, err)
	}
	return mongoClient.Database(dbName)
}
