package opnsense

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
			Status int `json:"status"`
		} `json:"System"`
		CrashReporter struct {
			Message    string `json:"message"`
			Status     string `json:"status"`
			StatusCode int    `json:"statusCode"`
		} `json:"CrashReporter"`
		Firewall struct {
			Message    string `json:"message"`
			Status     int    `json:"status"`
			StatusCode int    `json:"statusCode"`
		} `json:"Firewall"`
	} `json:"metadata"`
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
