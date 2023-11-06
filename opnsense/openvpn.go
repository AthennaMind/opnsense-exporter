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

func (c *Client) FetchOpenVPNInstances() (OpenVPNInstances, *APICallError) {
	var resp openVPNSearchResponse
	var data OpenVPNInstances

	url, ok := c.endpoints["openVPNInstances"]
	if !ok {
		return data, &APICallError{
			Endpoint:   "openvpn",
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
