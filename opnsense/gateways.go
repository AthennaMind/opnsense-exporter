package opnsense

import (
	"log/slog"
)

// GatewayStatus is the custom type that represents the status of a gateway
type GatewayStatusType int

const (
	GatewayStatusOffline GatewayStatusType = iota
	GatewayStatusOnline
	GatewayStatusUnknown
	GatewayStatusPeding
)

// gatewayConfigurationResponse is the response from the OPNsense API that contains the gateways configuration details
type gatewayConfigurationResponse struct {
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
	Rows     []struct {
		Disabled             bool   `json:"disabled"`
		Name                 string `json:"name"`
		Description          string `json:"descr"`
		HardwareInterface    string `json:"interface"`
		IPProtocol           string `json:"ipprotocol"`
		Gateway              string `json:"gateway"`
		DefaultGateway       bool   `json:"defaultgw"`
		FarGateway           string `json:"fargw"`
		MonitorDisable       string `json:"monitor_disable"`
		MonitorNoRoute       string `json:"monitor_noroute"`
		Monitor              string `json:"monitor"`
		ForceDown            string `json:"force_down"`
		Priority             string `json:"priority"`
		Weight               int    `json:"weight"`
		LatencyLow           string `json:"latencylow"`
		CurrentLatencyLow    string `json:"current_latencylow"`
		LatencyHigh          string `json:"latencyhigh"`
		CurrentLatencyHigh   string `json:"current_latencyhigh"`
		LossLow              string `json:"losslow"`
		CurrentLossLow       string `json:"current_losslow"`
		LossHigh             string `json:"losshigh"`
		CurrentLossHigh      string `json:"current_losshigh"`
		Interval             string `json:"interval"`
		CurrentInterval      string `json:"current_interval"`
		TimePeriod           string `json:"time_period"`
		CurrentTimePeriod    string `json:"current_time_period"`
		LossInterval         string `json:"loss_interval"`
		CurrentLossInterval  string `json:"current_loss_interval"`
		DataLength           string `json:"data_length"`
		CurrentDataLength    string `json:"current_data_length"`
		UUID                 string `json:"uuid"`
		Interface            string `json:"if"`
		Attribute            int    `json:"attribute"`
		Dynamic              bool   `json:"dynamic"`
		Virtual              bool   `json:"virtual"`
		Upstream             bool   `json:"upstream"`
		InterfaceDescription string `json:"interface_descr"`
		Status               string `json:"status"`
		Delay                string `json:"delay"`
		StdDev               string `json:"stddev"`
		Loss                 string `json:"loss"`
		LabelClass           string `json:"label_class"`
	} `json:"rows"`
}

type Gateway struct {
	Name                 string
	Description          string
	Enabled              bool
	HardwareInterface    string
	IPProtocol           string
	Gateway              string
	DefaultGateway       bool
	FarGateway           string
	MonitorEnabled       bool
	MonitorNoRoute       bool
	Monitor              string
	ForceDown            bool
	Priority             string
	Weight               int
	LatencyLow           string
	LatencyHigh          string
	LossLow              string
	LossHigh             string
	Interval             string
	TimePeriod           string
	LossInterval         string
	DataLength           string
	UUID                 string
	Interface            string
	Attribute            int
	Dynamic              bool
	Virtual              bool
	Upstream             bool
	InterfaceDescription string
	Status               GatewayStatusType
	Delay                float64
	StdDev               float64
	Loss                 float64
	LabelClass           string
}

type Gateways struct {
	Gateways []Gateway
}

// parseGatewayStatus parses a string status to a GatewayStatus type.
func parseGatewayStatus(statusTranslated string, logger *slog.Logger, originalStatus string) GatewayStatusType {
	switch statusTranslated {
	case "Online":
		return GatewayStatusOnline
	case "Offline":
		return GatewayStatusOffline
	case "Pending":
		return GatewayStatusPeding
	default:
		logger.Warn("unknown gateway status detected", "status", originalStatus)
		return GatewayStatusUnknown
	}
}

// FetchGateways fetches the gateways status details from the OPNsense API
// and returns a safe wrapper Gateways struct.
func (c *Client) FetchGateways() (Gateways, *APICallError) {
	var resp gatewayConfigurationResponse
	var data Gateways

	url, ok := c.endpoints["gatewaysStatus"]
	if !ok {
		return data, &APICallError{
			Endpoint:   "gateways",
			Message:    "endpoint not found in client endpoints",
			StatusCode: 0,
		}
	}
	err := c.do("GET", url, nil, &resp)
	if err != nil {
		return data, err
	}

	for _, v := range resp.Rows {
		delay := -1.0
		stdDev := -1.0
		loss := -1.0
		if !v.Disabled && !parseStringToBool(v.MonitorDisable) {
			delay = parseStringToFloatWithReplace(v.Delay, c.gatewayRTTRegex, " ms", "rtt", c.log)
			stdDev = parseStringToFloatWithReplace(v.StdDev, c.gatewayRTTRegex, " ms", "rttd", c.log)
			loss = parseStringToFloatWithReplace(v.Loss, c.gatewayLossRegex, " %", "loss", c.log)
		}

		g := Gateway{
			Name:                 v.Name,
			Description:          v.Description,
			Enabled:              !v.Disabled,
			HardwareInterface:    v.HardwareInterface,
			IPProtocol:           v.IPProtocol,
			Gateway:              v.Gateway,
			DefaultGateway:       v.DefaultGateway,
			FarGateway:           v.FarGateway,
			MonitorEnabled:       !parseStringToBool(v.MonitorDisable),
			MonitorNoRoute:       parseStringToBool(v.MonitorNoRoute),
			Monitor:              v.Monitor,
			ForceDown:            parseStringToBool(v.ForceDown),
			Priority:             v.Priority,
			Weight:               v.Weight,
			LatencyLow:           v.LatencyLow,
			LatencyHigh:          v.LatencyHigh,
			LossLow:              v.LossLow,
			LossHigh:             v.LossHigh,
			Interval:             v.Interval,
			TimePeriod:           v.TimePeriod,
			LossInterval:         v.LossInterval,
			DataLength:           v.DataLength,
			UUID:                 v.UUID,
			Interface:            v.Interface,
			Attribute:            v.Attribute,
			Dynamic:              v.Dynamic,
			Virtual:              v.Virtual,
			Upstream:             v.Upstream,
			InterfaceDescription: v.InterfaceDescription,
			Status:               parseGatewayStatus(v.Status, c.log, v.Status),
			Delay:                delay,
			StdDev:               stdDev,
			Loss:                 loss,
			LabelClass:           v.LabelClass,
		}

		switch {
		case g.LatencyLow == "":
			g.LatencyLow = v.CurrentLatencyLow
		case g.LatencyHigh == "":
			g.LatencyHigh = v.CurrentLatencyHigh
		case g.LossLow == "":
			g.LossLow = v.CurrentLossLow
		case g.LossHigh == "":
			g.LossHigh = v.CurrentLossHigh
		case g.Interval == "":
			g.Interval = v.CurrentInterval
		case g.TimePeriod == "":
			g.TimePeriod = v.CurrentTimePeriod
		case g.LossInterval == "":
			g.LossInterval = v.CurrentLossInterval
		case g.DataLength == "":
			g.DataLength = v.CurrentDataLength
		}

		data.Gateways = append(data.Gateways, g)
	}

	return data, nil
}
