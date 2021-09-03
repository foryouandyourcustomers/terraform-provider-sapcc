package client

import (
	"net/http"
	"terraform-provider-sapcc/internal/provider"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (c *Client) GetBuild(buildCode string) (*provider.Build, error) {
	request, err := http.NewRequest("GET", c.buildURL(buildCode), nil)
	if err != nil {
		return nil, err
	}

	resp, statusCode, err := c.doRequest(request)

	if err != nil {
		return nil, err
	}

	var build provider.Build

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
				var properties []provider.BuildProperty

				for _, v := range v.([]interface{}) {
					v := v.(map[string]interface{})

					properties = append(properties, provider.BuildProperty{
						Key: types.String{Value: v["key"].(string)},
						Val: types.String{Value: v["value"].(string)},
					})
				}

				build.Properties = properties

			default:
				logger.Debug("Unexpected data received from build response", hclog.Fmt(" k=%s v=%s, ignoring", k, v))
			}
		}
	}

	return &build, err
}
