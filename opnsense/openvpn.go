package opnsense

import "strings"

const fetchOpenVPNPayload = `{"current":1,"rowCount":-1,"sort":{},"searchPhrase":""}`

type openVPNSearchResponse struct {
	Rows []struct {
		UUID        string `json:"uuid"`
		Description string `json:"description"`
		Role        string `json:"role"`
		DevType     string `json:"dev_type"`
		Enabled     string `json:"enabled"`
	} `json:"rows"`
	RowCount int `json:"rowCount"`
	Total    int `json:"total"`
	Current  int `json:"current"`
}

type openVPNSearchSessionsResponse struct {
	Rows []struct {
		Description    string `json:"description"`
		Username       string `json:"username"`
		VirtualAddress string `json:"virtual_address"`
		Status         string `json:"status"`
	} `json:"rows"`
	RowCount int `json:"rowCount"`
	Total    int `json:"total"`
	Current  int `json:"current"`
}

type OpenVPN struct {
	UUID        string
	Description string
	Role        string
	DevType     string
	Enabled     int
}
type OpenVPNInstances struct {
	Rows []OpenVPN
}

type Sessions struct {
	Description    string
	Username       string
	VirtualAddress string
	Status         int
}
type OpenVPNSessions struct {
	Rows []Sessions
}

func (c *Client) FetchOpenVPNInstances() (OpenVPNInstances, *APICallError) {
	var resp openVPNSearchResponse
	var data OpenVPNInstances

	url, ok := c.endpoints["openVPNInstances"]
	if !ok {
		return data, &APICallError{
			Endpoint:   "openVPNInstances",
			Message:    "endpoint not found in client endpoints",
			StatusCode: 0,
		}
	}

	if err := c.do("POST", url, strings.NewReader(fetchOpenVPNPayload), &resp); err != nil {
		return data, err
	}

	for _, v := range resp.Rows {
		enabled, err := parseStringToInt(v.Enabled, url)
		if err != nil {
			return data, err
		}
		data.Rows = append(data.Rows, OpenVPN{
			UUID:        v.UUID,
			Description: v.Description,
			Role:        strings.ToLower(v.Role),
			DevType:     v.DevType,
			Enabled:     enabled,
		})
	}

	return data, nil
}

func (c *Client) FetchOpenVPNSessions() (OpenVPNSessions, *APICallError) {
	var resp openVPNSearchSessionsResponse
	var data OpenVPNSessions

	url, ok := c.endpoints["openVPNSessions"]
	if !ok {
		return data, &APICallError{
			Endpoint:   "openVPNSessions",
			Message:    "endpoint not found in client endpoints",
			StatusCode: 0,
		}
	}

	if err := c.do("GET", url, nil, &resp); err != nil {
		return data, err
	}

	for _, v := range resp.Rows {
		data.Rows = append(data.Rows, Sessions{
			Description:    v.Description,
			Username:       v.Username,
			VirtualAddress: v.VirtualAddress,
			Status:         parseOpenVPNsessionStatusToInt(v.Status),
		})
	}

	return data, nil
}
