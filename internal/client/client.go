package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	deployProgressURL func(deployCode string) string
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
		HTTPClient:        &http.Client{Timeout: 60 * time.Second},
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
		deployProgressURL: func(deployCode string) string {
			return fmt.Sprintf("%s/deployments/%s/progress", baseURL, deployCode)
		},
	}

	return &c, nil
}

func (c *Client) doRequest(request *http.Request) (map[string]interface{}, int, error) {
	request.Header = http.Header{
		"Authorization": []string{c.AuthToken},
		"Content-Type":  []string{"application/json"},
	}

	logger.Debug("Sending request", hclog.Fmt("%+v", request))

	res, err := c.HTTPClient.Do(request)

	logger.Debug("Raw response", hclog.Fmt("%+v", res))

	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	jsonResponse := make(map[string]interface{})

	if len(body) > 0 {
		err = json.Unmarshal(body, &jsonResponse)
		logger.Debug("Response ", hclog.Fmt(" %+v", jsonResponse), " statusCode: ", hclog.Fmt("%s", res.StatusCode))
	} else {
		logger.Warn("Response has no content length")
	}

	// let the client handle all the errors
	return jsonResponse, res.StatusCode, err
}
