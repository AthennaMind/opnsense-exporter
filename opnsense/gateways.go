package opnsense

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// gatewaysStatusResponse is the response from the OPNsense API that contains the gateways status details
// The data is constucted in this script:
// ---> https://github.com/opnsense/core/blob/master/src/opnsense/scripts/routes/gateway_status.php
// Following the reverse engineering of the call:
// ---> https://github.com/opnsense/core/blob/master/src/etc/inc/plugins.inc.d/dpinger.inc#L368
// From this file we know that Loss and Delay always have the same format of '%0.1f ms'
type gatewaysStatusResponse struct {
	Status string `json:"status"`
	Items  []struct {
		Name             string `json:"name"`
		Address          string `json:"address"`
		Status           string `json:"status"`
		Loss             string `json:"loss"`
		Delay            string `json:"delay"`
		Stddev           string `json:"stddev"`
		StatusTranslated string `json:"status_translated"`
	} `json:"items"`
}

// GatewayStatus is the custom type that represents the status of a gateway
type GatewayStatus int

const (
	GatewayStatusOffline GatewayStatus = iota
	GatewayStatusOnline
	GatewayStatusUnknown
)

type Gateway struct {
	Name             string
	Address          string
	Status           GatewayStatus
	RTTMilliseconds  float64
	RTTDMilliseconds float64
	LossPercentage   float64
}

type Gateways struct {
	Gateways []Gateway
}

// parseGatewayStatus parses a string status to a GatewayStatus type.
func parseGatewayStatus(statusTranslated string, logger log.Logger, originalStatus string) GatewayStatus {
	switch statusTranslated {
	case "Online":
		return GatewayStatusOnline
	case "Offline":
		return GatewayStatusOffline
	default:
		level.Warn(logger).
			Log("msg", "unknown gateway status detected", "status", originalStatus)
		return GatewayStatusUnknown
	}
}

// FetchGateways fetches the gateways status details from the OPNsense API
// and returns a safe wrapper Gateways struct.
func (c *Client) FetchGateways() (Gateways, *APICallError) {
	var resp gatewaysStatusResponse
	var data Gateways

	url, ok := c.endpoints["gatewaysStatus"]
	if !ok {
		return data, &APICallError{
			Endpoint:   "gatewaysStatus",
			Message:    "endpoint not found in client endpoints",
			StatusCode: 0,
		}
	}
	err := c.do("GET", url, nil, &resp)
	if err != nil {
		return data, err
	}

	for _, v := range resp.Items {
		data.Gateways = append(data.Gateways, Gateway{
			Name:             v.Name,
			Address:          v.Address,
			Status:           parseGatewayStatus(v.StatusTranslated, c.log, v.Status),
			RTTMilliseconds:  parseStringToFloatWithReplace(v.Delay, c.gatewayRTTRegex, " ms", "rtt", c.log),
			RTTDMilliseconds: parseStringToFloatWithReplace(v.Stddev, c.gatewayRTTRegex, " ms", "rttd", c.log),
			LossPercentage:   parseStringToFloatWithReplace(v.Loss, c.gatewayLossRegex, " %", "loss", c.log),
		})
	}

	return data, nil
}
