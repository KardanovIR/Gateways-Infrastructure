package server

import (
	"context"
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/models"
	"google.golang.org/grpc"
	"math/big"
	"sync"
	"testing"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth/logger"
)

const serverPort = "20001"

var oneServer sync.Once

func TestGrpcServerGasPrice(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	conn, err := startServerAndGetConnection(ctx)
	if err != nil {
		log.Error("connection to grpc server don't establish", err)
		t.Fail()
	}
	defer conn.Close()

	c := pb.NewCommonClient(conn)
	reply, err := c.GasPrice(ctx, &pb.GasPriceRequest{})
	if err != nil {
		log.Error("getting gas price fails", err)
		t.Fail()
	}
	if reply.GasPrice != "3" {
		t.Fail()
	}
}

func startServerAndGetConnection(ctx context.Context) (*grpc.ClientConn, error) {
	log := logger.FromContext(ctx)
	oneServer.Do(func() {
		go func() {
			err := InitAndStart(ctx, serverPort, &nodeClientMock{})
			if err != nil {
				log.Fatal("server can't start", err)
			}
		}()
	})
	return grpc.Dial(fmt.Sprint(":", serverPort), grpc.WithInsecure())
}

type nodeClientMock struct {
}

func (n *nodeClientMock) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return new(big.Int).SetInt64(3), nil
}

func (n *nodeClientMock) SuggestFee(ctx context.Context) (*big.Int, error) {
	return nil, nil
}

func (n *nodeClientMock) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	return nil, nil
}
func (n *nodeClientMock) GetNextNonce(ctx context.Context, address string) (uint64, error) {
	return 0, nil
}

func (n *nodeClientMock) GenerateAddress(ctx context.Context) (publicAddress string, err error) {
	return "", nil
}
func (n *nodeClientMock) IsAddressValid(ctx context.Context, address string) bool {
	return true
}

func (n *nodeClientMock) CreateRawTransaction(ctx context.Context, addressFrom string, addressTo string,
	amount *big.Int) ([]byte, error) {
	return nil, nil
}
func (n *nodeClientMock) SignTransaction(ctx context.Context, senderAddr string, rlpTx []byte) ([]byte, error) {
	return nil, nil
}
func (n *nodeClientMock) SendTransaction(ctx context.Context, rlpTx []byte) (txHash string, err error) {
	return "", nil
}
func (n *nodeClientMock) GetTxStatusByTxID(ctx context.Context, txID string) (models.TxStatus, error) {
	return "", nil
}
