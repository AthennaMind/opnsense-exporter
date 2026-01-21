package opnsense

import (
	"strconv"
)

type HealthCheckResponse struct {
	System struct {
		Status string `json:"status"`
	} `json:"System"`
	CrashReporter struct {
		Message    string `json:"message"`
		Status     string `json:"status"`
		StatusCode int    `json:"statusCode"`
	} `json:"CrashReporter"`
	Firewall struct {
		Message    string `json:"message"`
		Status     string `json:"status"`
		StatusCode int    `json:"statusCode"`
	} `json:"Firewall"`
	// OPNsense>25.1 has a different structure
	// See https://github.com/AthennaMind/opnsense-exporter/issues/48#issuecomment-2692494735
	Metadata struct {
		System struct {
			Status interface{} `json:"status"`
		} `json:"System"`
		CrashReporter struct {
			Message    string `json:"message"`
			Status     string `json:"status"`
			StatusCode int    `json:"statusCode"`
		} `json:"CrashReporter"`
		Firewall struct {
			Message    string      `json:"message"`
			Status     interface{} `json:"status"`
			StatusCode int         `json:"statusCode"`
		} `json:"Firewall"`
	} `json:"metadata"`
}

// GetMetadataSystemStatus converts the Status field to int, handling both string and int types
func (h *HealthCheckResponse) GetMetadataSystemStatus() int {
	switch v := h.Metadata.System.Status.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// GetMetadataFirewallStatus converts the Status field to int, handling both string and int types
func (h *HealthCheckResponse) GetMetadataFirewallStatus() int {
	switch v := h.Metadata.Firewall.Status.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

const (
	HealthCheckStatusOK = "OK"
	// OPNsense>25.1 has a different value
	// See https://github.com/AthennaMind/opnsense-exporter/issues/48#issuecomment-2692494735
	HealthCheckStatusOK_v25_1 = 2
)

// HealthCheck checks if the OPNsense is up and running.
func (c *Client) HealthCheck() (HealthCheckResponse, error) {
	var resp HealthCheckResponse

	path, ok := c.endpoints["healthCheck"]

	if !ok {
		return HealthCheckResponse{}, &APICallError{
			Endpoint:   "healthCheck",
			Message:    "endpoint not found",
			StatusCode: 0,
		}
	}

	if err := c.do("GET", path, nil, &resp); err != nil {
		return HealthCheckResponse{}, err
	}

	return resp, nil
}
