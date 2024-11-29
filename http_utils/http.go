package http_utils

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
	"sync"
	"time"
)

var (
	once         *sync.Once = new(sync.Once)
	httpInstance *httpClient
)

type HTTPClient interface {
	Post(url string, header map[string]string, body interface{}) (interface{}, error)
	Get(url string, header map[string]string) (string, error)
}

type httpClient struct {
	client *resty.Client
}

func GetHttpClient() HTTPClient {
	once.Do(
		func() {
			httpInstance = &httpClient{
				client: resty.New(),
			}

			httpInstance.client.SetTimeout(10 * time.Second)
		},
	)

	return httpInstance
}

func (c *httpClient) Post(url string, header map[string]string, body interface{}) (interface{}, error) {
	req := c.client.R()

	if len(header) > 0 {
		for key, value := range header {
			req.SetHeader(key, value)
		}
	}

	req.SetBody(body)

	resp, err := req.Post(url)
	if err != nil {
		return nil, fmt.Errorf("error during POST request: %v", err)
	}

	if resp.StatusCode() != 200 {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err != nil {
			return nil, fmt.Errorf("request failed with status code %d. Response: %s", resp.StatusCode(), resp.String())
		}

		if errorMessage, exists := errorResponse["message"]; exists {
			return nil, fmt.Errorf("ERROR: %v", errorMessage)
		}
		return nil, fmt.Errorf("request failed with status code %d. Response: %s", resp.StatusCode(), resp.String())
	}

	var successResponse map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &successResponse); err != nil {
		return nil, fmt.Errorf("error parsing success response: %v", err)
	}

	if data, exists := successResponse["data"]; exists {
		return data, nil
	}

	return successResponse, nil
}

func (c *httpClient) Get(url string, header map[string]string) (string, error) {
	req := c.client.R()

	if len(header) > 0 {
		for key, value := range header {
			req.SetHeader(key, value)
		}
	}

	resp, err := req.Get(url)
	if err != nil {
		return "", fmt.Errorf("error during GET request: %v", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("GET request failed with status code %d. Response: %s", resp.StatusCode(), resp.String())
	}

	return resp.String(), nil
}
