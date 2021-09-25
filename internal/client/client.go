package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
)

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
	clientUseragent   string
	logger            hclog.Logger
}

// NewClient -
func NewClient(providerVersion, baseURL, authToken string, logger hclog.Logger) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL can not be empty")
	}

	if authToken == "" {
		return nil, errors.New("authToken can not be empty")
	}

	c := Client{
		clientUseragent:   fmt.Sprintf("terraform-provider-sapcc/%s", providerVersion),
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
		logger: logger,
	}

	return &c, nil
}

func (c *Client) doRequest(request *http.Request) (map[string]interface{}, int, error) {
	request.Header = http.Header{
		"Authorization": []string{c.AuthToken},
		"Content-Type":  []string{"application/json"},
		"User-Agent":    []string{c.clientUseragent},
	}

	c.logger.Debug("Sending request", request)

	res, err := c.HTTPClient.Do(request)

	c.logger.Trace("Raw response", res)

	if err != nil {
		return nil, 0, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.logger.Error("Error: ", err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	jsonResponse := make(map[string]interface{})

	if len(body) > 0 {
		err = json.Unmarshal(body, &jsonResponse)
		c.logger.Trace("Response ", jsonResponse, " statusCode: ", res.StatusCode)

		if jsonResponse["title"] == "Api Exception" {
			err = fmt.Errorf("%s: %s", jsonResponse["message"], jsonResponse["detail"])

			return nil, res.StatusCode, err
		}
	} else {
		c.logger.Warn("Response has no content length")
	}

	// let the client handle all the errors
	return jsonResponse, res.StatusCode, err
}
