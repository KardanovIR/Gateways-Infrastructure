package listener_test

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/config"
	modelsBtc "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/models"
	repository2 "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/repository"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Btc/services"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/grpc/client"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/models"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/repositories"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/server"
	coreServices "github.com/wavesplatform/GatewaysInfrastructure/Listeners/Core/services"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

const (
	txHashForSearchingByTxHash = "97d3dc37d1382c01a182f13d2b56228b96546cae3150e75bf463e21cad125886"
	addressForSearching        = "2Mxd3wMiJEhHqcMPX8BrFwHxXSSsDvrrpJN"
)

var (
	callBackChannel = make(chan string, 2)
	initTestOnce    sync.Once
)

// block 1580086
// address1 -> clientAddress 0.02320000
// address1 -> clientAddress 0.002
// block 1580089
// clientAddress -> address1 0.005
// listen clientAddress
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
		if _, err := mongoDB.Collection(repositories.CChainState).DeleteOne(ctx, bson.D{{"chaintype", models.Btc}}); err != nil {
			log.Error(err)
		}
		if err := mongoDB.Collection(repository2.UnspentTxCollection).Drop(ctx); err != nil {
			log.Error(err)
		}
		if err := mongoDB.Collection(repositories.Ctasks).Drop(ctx); err != nil {
			log.Error(err)
		}
	}()

	// add tasks Transfer
	_, err = grpcClient.AddTask(ctx,
		&blockchain.AddTaskRequest{
			ListenTo:     &blockchain.ListenObject{Type: "Address", Value: addressForSearching},
			CallbackType: string(models.InitInTx),
			ProcessId:    "111111111",
			TaskType:     strconv.Itoa(int(models.Permanent))},
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

	// wait for receiving callback
	var InitInTxHashs = make([]string, 0)
	var finishProcess bool

	for i := 0; i < 4; i++ {
		select {
		case callback := <-callBackChannel:
			if callback == "FinishProcess" {
				finishProcess = true
			} else {
				InitInTxHashs = append(InitInTxHashs, callback)
			}
		case <-time.After(10 * time.Second):
			log.Error("so long waiting...")
			t.FailNow()
		}
	}
	assert.Equal(t, 3, len(InitInTxHashs))
	assert.Equal(t, "2bb3d8130ab1672e6d4a0816ec487df05f61d5459bff75aeadae0aa84f3e4ecd", InitInTxHashs[0])
	assert.Equal(t, "4abb78db4420cb7af2526ce62e361d8c25852214d9980483ef20ba5b5fb72e3e", InitInTxHashs[1])
	assert.Equal(t, "97d3dc37d1382c01a182f13d2b56228b96546cae3150e75bf463e21cad125886", InitInTxHashs[2])
	assert.True(t, finishProcess)
	services.GetNodeReader().Stop(ctx)

	// check inputs in mongo
	inputsCount, err := mongoDB.Collection(repository2.UnspentTxCollection).CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, int64(2), inputsCount)

	firstInput := mongoDB.Collection(repository2.UnspentTxCollection).FindOne(ctx, bson.D{{"txHash", "4abb78db4420cb7af2526ce62e361d8c25852214d9980483ef20ba5b5fb72e3e"}})
	var unspentTx1 modelsBtc.UnspentTx
	err = firstInput.Decode(&unspentTx1)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, "2Mxd3wMiJEhHqcMPX8BrFwHxXSSsDvrrpJN", unspentTx1.Address)
	assert.Equal(t, uint64(200000), unspentTx1.Amount)
	assert.Equal(t, uint32(0), unspentTx1.OutputN)

	secondInput := mongoDB.Collection(repository2.UnspentTxCollection).FindOne(ctx, bson.D{{"txHash", "97d3dc37d1382c01a182f13d2b56228b96546cae3150e75bf463e21cad125886"}})
	var unspentTx2 modelsBtc.UnspentTx
	err = secondInput.Decode(&unspentTx2)
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
	assert.Equal(t, "2Mxd3wMiJEhHqcMPX8BrFwHxXSSsDvrrpJN", unspentTx2.Address)
	assert.Equal(t, uint64(1803520), unspentTx2.Amount)
	assert.Equal(t, uint32(1), unspentTx2.OutputN)

}

func TestNodeClient_GetBlockVerboseTx(t *testing.T) {
	ctx := context.Background()
	beforeTests(ctx, t)
	b, err := services.GetNodeClient().GetBlockVerboseTx(ctx, "000000006ba88d4ae2349a2e286752d649e0f18abf2eec99fce3c5615f1c02d8")
	assert.Nil(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, "000000006ba88d4ae2349a2e286752d649e0f18abf2eec99fce3c5615f1c02d8", b.Hash)
	assert.Equal(t, 374, len(b.Tx))
	assert.Equal(t, "4271e6e08b5f66e73b0d7a637f91ff70500b091f81027eca38c0459ff1a3e1b5", b.Tx[1].Txid)
	assert.Equal(t, 2, len(b.Tx[1].Vout))
	assert.Equal(t, uint32(1), b.Tx[1].Vout[1].N)
	assert.Equal(t, float64(0.76816512), b.Tx[1].Vout[1].Value)
	amount, err := services.GetIntFromFloat(ctx, b.Tx[1].Vout[1].Value)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, uint64(76816512), amount)
	assert.Equal(t, 1, len(b.Tx[1].Vin))
	assert.Equal(t, uint32(1), b.Tx[1].Vin[0].Vout)
	assert.Equal(t, "1beb61055b1b1e1f91ebff199ba780bdea08afa59bf528c59d20844cc49c4f25", b.Tx[1].Vin[0].Txid)

}

func beforeTests(ctx context.Context, t *testing.T) {
	log, _ := logger.Init(false, logger.DEBUG)
	initTestOnce.Do(func() {
		err := config.Load("./testdata/config_test.yml")
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("Cfg: %+v", config.Cfg)
		repository, err := repository2.New(ctx, config.Cfg.Db)
		if err != nil {
			log.Fatal(err)
		}

		if err := coreServices.NewCallbackService(ctx, config.Cfg.CallbackUrl, config.Cfg.Node.ChainType); err != nil {
			log.Fatal("Can't create callback service: ", err)
		}

		nodeClient := services.NewNodeClient(ctx, config.Cfg.Node)
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
	assert.Equal(c.t, addressForSearching, in.Address)
	callBackChannel <- in.TxHash
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
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+url))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB at %s: %v", url, err)
	}
	return mongoClient.Database(dbName)
}
