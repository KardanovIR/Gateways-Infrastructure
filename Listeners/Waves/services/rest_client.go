package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Waves/models"
)

type IRestClient interface {
	RequestCallback(ctx context.Context, callback models.Callback, params map[string]interface{}) error
	//Start(ctx context.Context) (err error)
	//Stop(ctx context.Context)
}

type restClient struct {
	//nodeClient    *ethclient.Client
	//rp            repositories.IRepository
	//conf          *config.Node
	//stopListenBTC chan struct{}
}

var (
	rc             IRestClient
	onceRestClient sync.Once
)

// New create node's client with connection to Waves node
func NewRestClient(ctx context.Context) error {
	log := logger.FromContext(ctx)
	var err error
	onceRestClient.Do(func() {
		rc = &restClient{}
	})

	if err != nil {
		log.Errorf("error during initialise rest client: %s", err)
		return err
	}

	return nil
}

// GetNodeReader returns node's reader instance.
// Client must be previously created with New(), in another case function throws panic
func GetRestClient() IRestClient {
	onceRestClient.Do(func() {
		panic("try to get rest client before it's creation!")
	})
	return rc
}

func (rc *restClient) RequestCallback(ctx context.Context, callback models.Callback, params map[string]interface{}) error {
	log := logger.FromContext(ctx)

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		log.Errorf("Error: %s", err)
		return err
	}

	_, err = rc.request(ctx, callback.Url, string(callback.Type), bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Errorf("Error: %s", err)
		return err
	}

	return nil
}

func (rc *restClient) request(ctx context.Context, urlPath, requestType string, data *bytes.Buffer) ([]byte, error) {
	log := logger.FromContext(ctx)

	req, err := http.NewRequest(requestType, urlPath, data)
	if err != nil {
		log.Errorf("Error: %s", err)
		return nil, err
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		log.Errorf("Error: %s", err)
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request to %s returned wrong status code $s", urlPath, res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error: %s", err)
		return nil, err
	}

	return body, nil
}
