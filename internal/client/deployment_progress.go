package client

import (
	"math/big"
	"net/http"
	"terraform-provider-sapcc/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (c *Client) GetDeploymentProgress(deploymentCode string) (*models.DeployProgress, int, error) {
	request, err := http.NewRequest("GET", c.deployProgressURL(deploymentCode), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, statusCode, err := c.doRequest(request)

	if err != nil {
		return nil, statusCode, err
	}

	var progress models.DeployProgress

	if statusCode == 200 {
		for k, v := range resp {
			switch k {
			case "subscriptionCode":
				progress.SubscriptionCode = types.String{Value: v.(string)}
			case "deploymentCode":
				progress.DeployCode = types.String{Value: v.(string)}
			case "deploymentStatus":
				progress.DeployStatus = types.String{Value: v.(string)}
			case "percentage":
				progress.ProgressPercentage = types.Number{Value: big.NewFloat(v.(float64))}
			default:
				c.logger.With("deploymentCode", deploymentCode).Debug("skipping {", k, v, "}")
			}
		}
	}

	return &progress, statusCode, err
}
