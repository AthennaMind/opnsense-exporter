package opnsense

type HealthCheckResponse struct {
	CrashReporter struct {
		StatusCode  int    `json:"statusCode"`
		Message     string `json:"message"`
		LogLocation string `json:"logLocation"`
		Timestamp   string `json:"timestamp"`
		Status      string `json:"status"`
	} `json:"CrashReporter"`
	Firewall struct {
		StatusCode  int    `json:"statusCode"`
		Message     string `json:"message"`
		LogLocation string `json:"logLocation"`
		Timestamp   string `json:"timestamp"`
		Status      string `json:"status"`
	} `json:"Firewall"`
	System struct {
		Status string `json:"status"`
	} `json:"System"`
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
