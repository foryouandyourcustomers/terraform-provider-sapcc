package client

import (
	"net/http"
	"terraform-provider-sapcc/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (c *Client) GetBuild(buildCode string) (*models.Build, int, error) {
	request, err := http.NewRequest("GET", c.buildURL(buildCode), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, statusCode, err := c.doRequest(request)

	if err != nil {
		return nil, statusCode, err
	}

	var build models.Build

	if statusCode == 200 {
		for k, v := range resp {
			switch k {
			case "createdBy":
				build.CreatedBy = types.String{Value: v.(string)}
			case "branch":
				build.Branch = types.String{Value: v.(string)}
			case "applicationDefinitionVersion":
				build.ApplicationDefinitionVersion = types.String{Value: v.(string)}
			case "subscriptionCode":
				build.SubscriptionCode = types.String{Value: v.(string)}
			case "applicationCode":
				build.ApplicationCode = types.String{Value: v.(string)}
			case "code":
				build.Code = types.String{Value: v.(string)}
			case "buildVersion":
				build.BuildVersion = types.String{Value: v.(string)}
			case "name":
				build.Name = types.String{Value: v.(string)}
			case "buildStartTimestamp":
				build.BuildStartTimestamp = types.String{Value: v.(string)}
			case "buildEndTimestamp":
				build.BuildEndTimestamp = types.String{Value: v.(string)}
			case "status":
				build.Status = types.String{Value: v.(string)}
			case "properties":
				var properties []models.BuildProperty

				for _, v := range v.([]interface{}) {
					v := v.(map[string]interface{})

					properties = append(properties, models.BuildProperty{
						Key: types.String{Value: v["key"].(string)},
						Val: types.String{Value: v["value"].(string)},
					})
				}

				build.Properties = properties

			default:
				c.logger.With("buildCode", buildCode).Debug("Unexpected data received from build response {", k, v, "}")
			}
		}
	}

	return &build, statusCode, err
}
