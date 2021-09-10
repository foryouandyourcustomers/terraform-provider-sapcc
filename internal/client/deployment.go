package client

import (
	"bytes"
	"fmt"
	"net/http"
	"terraform-provider-sapcc/internal/models"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (c *Client) GetDeployment(deploymentCode string) (*models.Deployment, int, error) {
	request, err := http.NewRequest("GET", c.deployURL(deploymentCode), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, statusCode, err := c.doRequest(request)

	if err != nil {
		return nil, statusCode, err
	}

	var deployment models.Deployment

	if statusCode == 200 {
		for k, v := range resp {
			var val string
			// we get many nil elements
			if v != nil {
				stringVal, ok := v.(string)
				if ok {
					val = stringVal
				} else {
					logger.Error("Nonstring type received", hclog.Fmt("key:%s value:%s", k, v))
				}
			}

			switch k {
			case "createdBy":
				deployment.CreatedBy = types.String{Value: val}
			case "buildCode":
				deployment.BuildCode = types.String{Value: val}
			case "createdTimestamp":
				deployment.CreatedTimestamp = types.String{Value: val}
			case "subscriptionCode":
				deployment.SubscriptionCode = types.String{Value: val}
			case "environmentCode":
				deployment.EnvironmentCode = types.String{Value: val}
			case "code":
				deployment.Code = types.String{Value: val}
			case "databaseUpdateMode":
				deployment.DatabaseUpdateMode = types.String{Value: val}
			case "strategy":
				deployment.Strategy = types.String{Value: val}
			case "scheduledTimestamp":
				deployment.ScheduledTimestamp = types.String{Value: val}
			case "deployedTimestamp":
				deployment.DeployedTimestamp = types.String{Value: val}
			case "undeployedTimestamp":
				deployment.UndeployedTimestamp = types.String{Value: val}
			case "failedTimestamp":
				deployment.FailedTimestamp = types.String{Value: val}
			case "status":
				deployment.Status = types.String{Value: val}
			case "cancelation":
				var cancelation []models.DeployCancellation

				if v != nil {
					v := v.(map[string]interface{})

					cancelation = append(cancelation, models.DeployCancellation{
						CancelledBy:      types.String{Value: v["canceledBy"].(string)},
						StartTimestamp:   types.String{Value: v["startTimestamp"].(string)},
						FinishTimestamp:  types.String{Value: v["finishedTimestamp"].(string)},
						Failed:           types.Bool{Value: v["failed"].(bool)},
						RollbackDatabase: types.Bool{Value: v["rollbackDatabase"].(bool)},
					})
				} else {
					// temporary fix
					cancelation = []models.DeployCancellation{}
				}

				deployment.Cancelation = cancelation
			default:
				logger.Debug("Unexpected data received from build response", hclog.Fmt(" k=%s v=%s, ignoring", k, v))
			}
		}
	}

	return &deployment, statusCode, err
}

func (c *Client) CreateDeployment(plan *models.Deployment) (*models.Deployment, int, error) {
	// prepare the deployment request
	var deployReq = []byte(fmt.Sprintf(`{"buildCode": "%s","databaseUpdateMode": "%s","environmentCode": "%s","strategy": "%s"}`,
		plan.BuildCode.Value,
		plan.DatabaseUpdateMode.Value,
		plan.EnvironmentCode.Value,
		plan.Strategy.Value))

	request, err := http.NewRequest("POST", c.DeploymentBaseURL, bytes.NewBuffer(deployReq))

	if err != nil {
		return nil, 0, err
	}

	resp, statusCode, err := c.doRequest(request)

	if err != nil {
		return nil, statusCode, err
	}

	if statusCode == 200 || statusCode == 201 {
		deploymentCode, ok := resp["code"].(string)

		if !ok {
			logger.Error("Unexpected data received, expected 'code' to be string: response", hclog.Fmt("response %s", resp))
			return nil, statusCode, err
		}

		return &models.Deployment{
			Code:               types.String{Value: deploymentCode},
			BuildCode:          types.String{Value: plan.BuildCode.Value},
			DatabaseUpdateMode: types.String{Value: plan.DatabaseUpdateMode.Value},
			EnvironmentCode:    types.String{Value: plan.EnvironmentCode.Value},
			Strategy:           types.String{Value: plan.Strategy.Value},
		}, statusCode, err
	}

	return nil, statusCode, err
}
