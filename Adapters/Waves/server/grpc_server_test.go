package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"sync"
	"testing"

	pb "github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/grpc"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/models"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

const serverPort = "20001"

var oneServer sync.Once

func TestGrpcServerGetLastBlockHeight(t *testing.T) {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	conn, err := startServerAndGetConnection(ctx)
	if err != nil {
		log.Error("connection to grpc server don't establish", err)
		t.Fail()
	}
	defer conn.Close()

	c := pb.NewCommonClient(conn)
	reply, err := c.GetLastBlockHeight(ctx, &pb.BlockRequest{})
	if err != nil {
		log.Error("getting last block height fails", err)
		t.Fail()
	}
	if reply.Block != "3" {
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

func (n *nodeClientMock) GetLastBlockHeight(ctx context.Context) (string, error) {
	return "3", nil
}

func (n *nodeClientMock) GenerateAddress(ctx context.Context) (publicAddress string, err error) {
	return "", nil
}

func (n *nodeClientMock) ValidateAddress(ctx context.Context, address string) (bool, error) {
	return true, nil
}

func (n *nodeClientMock) GetBalance(ctx context.Context, address string) (uint64, error) {
	return 0, nil
}
func (n *nodeClientMock) GetAllBalances(ctx context.Context, address string) (*models.AccountBalance, error) {
	return nil, nil
}

func (n *nodeClientMock) Fee(ctx context.Context, senderPublicKey string, feeAssetId string) (uint64, error) {
	return 0, nil
}

func (n *nodeClientMock) FeeForTx(ctx context.Context, tx *proto.TransferV2) (uint64, error) {
	return 0, nil
}

func (n *nodeClientMock) CreateRawTxBySendersAddress(ctx context.Context, addressFrom string, addressTo string, amount uint64) ([]byte, error) {
	return nil, nil
}

func (n *nodeClientMock) CreateRawTxBySendersPublicKey(ctx context.Context, sendersPublicKey string, addressTo string, amount uint64, assetId string) ([]byte, error) {
	return nil, nil
}

func (n *nodeClientMock) SignTxWithKeepedSecretKey(ctx context.Context, sendersAddress string, txUnsigned []byte) ([]byte, error) {
	return nil, nil
}

func (n *nodeClientMock) SignTxWithSecretKey(ctx context.Context, secretKeyInBase58 string, txUnsigned []byte) ([]byte, error) {
	return nil, nil
}

func (n *nodeClientMock) SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error) {
	return "", nil
}

func (n *nodeClientMock) GetTransactionByTxId(ctx context.Context, txId string) ([]byte, error) {
	return nil, nil
}

func (n *nodeClientMock) GetTransactionStatus(ctx context.Context, txId string) (models.TxStatus, error) {
	return models.TxStatusUnKnown, nil
}
