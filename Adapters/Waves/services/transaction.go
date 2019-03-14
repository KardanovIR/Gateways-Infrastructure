package services

import (
	"context"
	"fmt"
	"time"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

// CreateRawTxBySendersAddress creates transaction for senders address if private key keeps in adapter
func (cl *nodeClient) CreateRawTxBySendersAddress(ctx context.Context, addressFrom string,
	addressTo string, amount uint64) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'CreateRawTxBySendersAddress' from %s to %s amount %d",
		addressFrom, addressTo, amount)
	secretKey, ok := cl.privateKeys[addressFrom]
	if !ok {
		return nil, fmt.Errorf("haven't private key for address %s", addressFrom)
	}
	senderPublic := crypto.GeneratePublicKey(secretKey)
	return cl.createRawTransaction(ctx, senderPublic, addressTo, amount)
}

// CreateRawTxBySendersPublicKey creates transaction using public key. Private key is not used
func (cl *nodeClient) CreateRawTxBySendersPublicKey(ctx context.Context, sendersPublicKey string,
	addressTo string, amount uint64) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'CreateRawTxBySendersPublicKey' pk %s send amount %d to %s",
		sendersPublicKey, amount, addressTo)

	senderPublic, err := crypto.NewPublicKeyFromBase58(sendersPublicKey)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return cl.createRawTransaction(ctx, senderPublic, addressTo, amount)
}

func (cl *nodeClient) createRawTransaction(ctx context.Context, senderPublic crypto.PublicKey,
	addressTo string, amount uint64) ([]byte, error) {
	log := logger.FromContext(ctx)
	amountAsset := proto.OptionalAsset{}
	feeAsset := proto.OptionalAsset{}
	recipientAddress, err := proto.NewAddressFromString(addressTo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	fee, err := cl.Fee(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	timestamp := time.Now().Unix() * 1000
	tx, err := proto.NewUnsignedTransferV2(senderPublic, amountAsset, feeAsset, uint64(timestamp), amount, fee,
		recipientAddress, "")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	txBinary, err := tx.BodyMarshalBinary()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return txBinary, nil
}

func (cl *nodeClient) SignTxWithKeepedSecretKey(ctx context.Context, sendersAddress string, txUnsigned []byte) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'SignTxWithKeepedSecretKey' for address %s", sendersAddress)
	secretKey, ok := cl.privateKeys[sendersAddress]
	if !ok {
		return nil, fmt.Errorf("haven't private key for address %s", sendersAddress)
	}
	return cl.signTransaction(ctx, secretKey, txUnsigned)
}

func (cl *nodeClient) SignTxWithSecretKey(ctx context.Context, secretKeyInBase58 string, txUnsigned []byte) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'SignTxWithSecretKey'")
	secretKey, err := crypto.NewSecretKeyFromBase58(secretKeyInBase58)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return cl.signTransaction(ctx, secretKey, txUnsigned)
}

func (cl *nodeClient) signTransaction(ctx context.Context, secretKey crypto.SecretKey, txUnsigned []byte) ([]byte, error) {
	log := logger.FromContext(ctx)
	tx := new(proto.TransferV2)
	if err := tx.BodyUnmarshalBinary(txUnsigned); err != nil {
		log.Error(err)
		return nil, err
	}
	if err := tx.Sign(secretKey); err != nil {
		log.Error(err)
		return nil, err
	}
	txBinary, err := tx.MarshalBinary()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return txBinary, nil
}

func (cl *nodeClient) SendTransaction(ctx context.Context, txSigned []byte) (txId string, err error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'SendTransaction'")
	tx := new(proto.TransferV2)
	if err := tx.UnmarshalBinary(txSigned); err != nil {
		log.Error(err)
		return "", err
	}
	if tx.ID != nil {
		txId = tx.ID.String()
	}
	log.Debugf("try to send tx with ID %s", tx.ID)
	_, err = cl.nodeClient.Transactions.Broadcast(ctx, tx)
	if err != nil {
		log.Error("sending tx fails", err)
	}
	return
}

func (cl *nodeClient) GetTransactionByTx(ctx context.Context, txId string) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'GetTransactionByTx' for txID %s", txId)

	id, err := crypto.NewDigestFromBase58(txId)
	tr, _, err := cl.nodeClient.Transactions.Info(ctx, id)
	if err != nil {
		log.Error("getting tx fails", err)
	}
	b, err := tr.MarshalBinary()
	if err != nil {
		log.Error("can't marshall binary", err)
	}
	return b, nil
}
