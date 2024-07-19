package opnsense

import "fmt"

// APICallError is an error returned by the OPNsense API
type APICallError struct {
	Endpoint   string
	Message    string
	StatusCode int
}

func (e APICallError) Error() string {
	return fmt.Sprintf(
		"opnsense-client api call error: endpoint: %s; failed status code: %d; msg: %s", e.Endpoint, e.StatusCode, e.Message,
	)
}
