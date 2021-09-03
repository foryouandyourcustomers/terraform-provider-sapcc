package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
)

var logger = hclog.New(&hclog.LoggerOptions{
	Name:  "sapcc-http-client",
	Level: hclog.Debug,
})

// Client -
type Client struct {
	HTTPClient        *http.Client
	BaseURL           string
	DeploymentBaseURL string
	BuildsBaseURL     string
	AuthToken         string
	buildURL          func(buildCode string) string
	deployURL         func(deployCode string) string
}

// NewClient -
func NewClient(baseURL, authToken string) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL can not be empty")
	}

	if authToken == "" {
		return nil, errors.New("authToken can not be empty")
	}

	c := Client{
		HTTPClient:        &http.Client{Timeout: 10 * time.Second},
		AuthToken:         authToken,
		BaseURL:           baseURL,
		BuildsBaseURL:     fmt.Sprintf("%s/builds", baseURL),
		DeploymentBaseURL: fmt.Sprintf("%s/deployments", baseURL),
		buildURL: func(buildCode string) string {
			return fmt.Sprintf("%s/builds/%s", baseURL, buildCode)
		},
		deployURL: func(deployCode string) string {
			return fmt.Sprintf("%s/deployments/%s", baseURL, deployCode)
		},
	}

	return &c, nil
}

func (c *Client) doRequest(request *http.Request) (map[string]interface{}, int, error) {
	request.Header = http.Header{
		"Authorization": []string{c.AuthToken},
		"Content-Type":  []string{"application/json"},
	}

	res, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, res.StatusCode, err
	}
	defer res.Body.Close()

	jsonResponse := make(map[string]interface{})

	err = json.NewDecoder(res.Body).Decode(&jsonResponse)
	// let the client handle all the errors
	return jsonResponse, res.StatusCode, err
}
