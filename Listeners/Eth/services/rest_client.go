package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/logger"
	"github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth/models"
	"io/ioutil"
	"net/http"
	"sync"
)

type IRestClient interface {
	RequestCallback(ctx context.Context, callback models.Callback, params map[string]interface{}) (tasks map[string]interface{}, err error)
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

// New create node's client with connection to eth node
func NewRestClient(ctx context.Context) error {
	log := logger.FromContext(ctx)
	var err error
	onceRestClient.Do(func() {
		rc = &restClient{}
	})

	if err != nil {
		log.Errorf("error during initialise node client: %s", err)
		return err
	}

	return nil
}

// GetNodeReader returns node's reader instance.
// Client must be previously created with New(), in another case function throws panic
func GetRestClient() IRestClient {
	onceRestClient.Do(func() {
		panic("try to get node reader before it's creation!")
	})
	return rc
}

func (rc *restClient) RequestCallback(ctx context.Context, callback models.Callback, params map[string]interface{}) (tasks map[string]interface{}, err error) {
	log := logger.FromContext(ctx)
	var response map[string]interface{}

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		log.Errorf("Error: %s", err)
		return response, err
	}

	response, err = rc.request(ctx, callback.Url, string(callback.Type), "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Errorf("Error: %s", err)
		return response, err
	}

	if response["errors"] != nil {
		errorsList, err := json.Marshal(response["errors"])

		if err != nil {
			log.Errorf("Error: %s", err)
			return response, err
		}
		log.Errorf("Errors: %s", errorsList)
		return response, errors.New(string(errorsList))
	}

	return response, nil
}

func (rc *restClient) request(ctx context.Context, urlPath, requestType, contentType string, data *bytes.Buffer) (map[string]interface{}, error) {
	log := logger.FromContext(ctx)
	var response map[string]interface{}

	req, err := http.NewRequest(requestType, urlPath, data)
	if err != nil {
		log.Errorf("Error: %s", err)
		return response, err
	}

	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		log.Errorf("Error: %s", err)
		return response, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("Request to %s returned wrong status code", urlPath)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error: %s", err)
		return response, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Errorf("Error: %s", err)
		return response, err
	}

	return response, nil
}
