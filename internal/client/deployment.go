package client

import (
	"net/http"
	"terraform-provider-sapcc/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/go-hclog"
)

func (c *Client) GetDeployment(deploymentCode string) (*provider.Deployment, error) {
	request, err := http.NewRequest("GET", c.deployURL(deploymentCode), nil)
	if err != nil {
		return nil, err
	}

	resp, statusCode, err := c.doRequest(request)

	if err != nil {
		return nil, err
	}

	var deployment provider.Deployment

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
				var cancelation []provider.DeployCancellation

				if v != nil {
					v := v.(map[string]interface{})

					cancelation = append(cancelation, provider.DeployCancellation{
						CancelledBy:      types.String{Value: v["canceledBy"].(string)},
						StartTimestamp:   types.String{Value: v["startTimestamp"].(string)},
						FinishTimestamp:  types.String{Value: v["finishedTimestamp"].(string)},
						Failed:           types.Bool{Value: v["failed"].(bool)},
						RollbackDatabase: types.Bool{Value: v["rollbackDatabase"].(bool)},
					})
				}

				deployment.Cancelation = cancelation
			default:
				logger.Debug("Unexpected data received from build response", hclog.Fmt(" k=%s v=%s, ignoring", k, v))
			}
		}
	}

	return &deployment, err
}
