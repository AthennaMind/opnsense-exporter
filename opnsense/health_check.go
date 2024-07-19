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
}

const HealthCheckStatusOK = "OK"

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
