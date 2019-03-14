package listener_test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/config"
	pb "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/repositories"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/server"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/services"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
	"testing"
	"time"
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

	// add tasks Transfer
	_, err = grpcClient.AddTask(ctx,
		&pb.AddTaskRequest{
			Address: "0x9515735d60e8ff4036efaffaf3370f3097615d19", CallbackType: string(models.Get),
			CallbackUrl: fmt.Sprintf("http://localhost:%s/transfer", httpServerPort),
			TaskType:    strconv.Itoa(int(models.OneTime))})

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
			if callback == "transfer" {
				isTransfer = true
			}
		case <-time.After(10 * time.Second):
			log.Error("so long waiting...")
			t.FailNow()
		}

	if !isTransfer {
		t.Fail()
	}



	defer func() {
		err := mongoClient.DropDatabase()
		if err != nil {
			log.Error("node reader task fails", err)
		}
	}()

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
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func mongoConnect(ctx context.Context, url string, dbName string) *mgo.Database {
	log := logger.FromContext(ctx)
	conn, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB at %s: %v", url, err)
	}
	return conn.DB(dbName)
}
