package services

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type WrongCodeError struct {
	Code int
	Body string
}

func (e WrongCodeError) Error() string {
	return fmt.Sprintf("wrong http status code %d, body %s", e.Code, e.Body)
}

func (cl *nodeClient) Request(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, &WrongCodeError{Code: resp.StatusCode, Body: string(body)}

	}
	return ioutil.ReadAll(resp.Body)
}
