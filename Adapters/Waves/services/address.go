package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"

	"github.com/wavesplatform/GatewaysInfrastructure/Adapters/Waves/logger"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func (cl *nodeClient) GenerateAddress(ctx context.Context) (publicAddress string, err error) {
	log := logger.FromContext(ctx)
	log.Info("call service method 'GenerateAddress'")
	seed := make([]byte, crypto.SecretKeySize)
	_, err = io.ReadFull(rand.Reader, seed)
	if err != nil {
		return
	}
	secret, public := crypto.GenerateKeyPair(seed)
	address, err := proto.NewAddressFromPublicKey(cl.chainID.Schema(), public)
	if err != nil {
		return
	}
	log.Debugf("privateKey: %s; publicKey: %s; address: %s", secret.String(), public.String(), address)
	cl.privateKeys[address.String()] = secret
	return address.String(), nil
}

func (cl *nodeClient) ValidateAddress(ctx context.Context, address string) (bool, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'ValidateAddress' for %s", address)
	a, err := proto.NewAddressFromString(address)
	if err != nil {
		return false, err
	}
	if a[1] != cl.chainID.Schema() {
		return false, fmt.Errorf("address for network %s, client connected with %s",
			string(a[1]), string(cl.chainID.Schema()))
	}
	return a.Valid()
}

func (cl *nodeClient) CheckAddress(ctx context.Context, address string, assetID string) (bool, error) {
	log := logger.FromContext(ctx)
	log.Infof("call service method 'CheckAddress' for %s", address)
	if len(cl.CheckAddressUrl) == 0 {
		return true, nil
	}
	url := fmt.Sprintf("%s/%s", cl.CheckAddressUrl, address)
	log.Infof("call external service for check address for %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("X-API-Key", cl.nodeClient.GetOptions().ApiKey)
	assetsList := make([]string, 0)
	if _, err := cl.nodeClient.Do(ctx, req, &assetsList); err != nil {
		log.Errorf("CheckAddress fails: %s", err)
		return false, err
	}
	for _, a := range assetsList {
		if a == assetID {
			return false, nil
		}
	}
	return true, nil
}
