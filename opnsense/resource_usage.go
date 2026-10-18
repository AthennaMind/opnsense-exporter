package opnsense

import (
	"encoding/json"
	"strconv"
)

// ResourceUsage holds current count and capacity reported by OPNsense.
type ResourceUsage struct {
	Current uint64
	Limit   uint64
}

func (c *Client) FetchPFStates() (ResourceUsage, *APICallError) {
	var resp struct {
		Current json.Number `json:"current"`
		Limit   json.Number `json:"limit"`
	}

	path, ok := c.endpoints["pfStates"]
	if !ok {
		return ResourceUsage{}, &APICallError{Endpoint: "pfStates", Message: "endpoint not found"}
	}
	if err := c.do("GET", path, nil, &resp); err != nil {
		return ResourceUsage{}, err
	}
	return parseResourceUsage(resp.Current, resp.Limit, path)
}

func (c *Client) FetchMbufUsage() (ResourceUsage, *APICallError) {
	var resp struct {
		Statistics struct {
			Current json.Number `json:"cluster-total"`
			Limit   json.Number `json:"cluster-max"`
		} `json:"mbuf-statistics"`
	}

	path, ok := c.endpoints["systemMbuf"]
	if !ok {
		return ResourceUsage{}, &APICallError{Endpoint: "systemMbuf", Message: "endpoint not found"}
	}
	if err := c.do("GET", path, nil, &resp); err != nil {
		return ResourceUsage{}, err
	}
	return parseResourceUsage(resp.Statistics.Current, resp.Statistics.Limit, path)
}

func parseResourceUsage(current, limit json.Number, path EndpointPath) (ResourceUsage, *APICallError) {
	// Missing fields and backend failure responses must not become zero usage.
	currentValue, err := strconv.ParseUint(current.String(), 10, 64)
	if err != nil {
		return ResourceUsage{}, &APICallError{
			Endpoint: string(path),
			Message:  "invalid current resource count: " + err.Error(),
		}
	}
	limitValue, err := strconv.ParseUint(limit.String(), 10, 64)
	if err != nil {
		return ResourceUsage{}, &APICallError{
			Endpoint: string(path),
			Message:  "invalid resource limit: " + err.Error(),
		}
	}
	return ResourceUsage{Current: currentValue, Limit: limitValue}, nil
}
