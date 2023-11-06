package opnsense

import (
	"strings"
)

type arpSearchResponse struct {
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
	Rows     []struct {
		Mac             string `json:"mac"`
		IP              string `json:"ip"`
		Intf            string `json:"intf"`
		Expired         bool   `json:"expired"`
		Expires         int    `json:"expires"`
		Permanent       bool   `json:"permanent"`
		Type            string `json:"type"`
		Manufacturer    string `json:"manufacturer"`
		Hostname        string `json:"hostname"`
		IntfDescription string `json:"intf_description"`
	} `json:"rows"`
}

type Arp struct {
	Mac             string
	IP              string
	Expired         bool
	Expires         int
	Permanent       bool
	Type            string
	Hostname        string
	IntfDescription string
}

type ArpTable struct {
	TotalEntries int
	Arp          []Arp
}

const fetchArpPayload = `{"current":1,"rowCount":-1,"sort":{},"searchPhrase":"","resolve":"no"}`

func (c *Client) FetchArpTable() (ArpTable, *APICallError) {
	var resp arpSearchResponse
	var arpTable ArpTable

	path, ok := c.endpoints["arp"]
	if !ok {
		return arpTable, &APICallError{
			Endpoint:   "arp",
			Message:    "endpoint not found",
			StatusCode: 0,
		}
	}

	if err := c.do("POST", path, strings.NewReader(fetchArpPayload), &resp); err != nil {
		return arpTable, err
	}

	for _, arp := range resp.Rows {
		a := Arp{
			Mac:             arp.Mac,
			IP:              arp.IP,
			Expired:         arp.Expired,
			Expires:         arp.Expires,
			Permanent:       arp.Permanent,
			Type:            arp.Type,
			Hostname:        arp.Hostname,
			IntfDescription: arp.IntfDescription,
		}
		arpTable.Arp = append(arpTable.Arp, a)
	}

	arpTable.TotalEntries = resp.Total

	return arpTable, nil
}
