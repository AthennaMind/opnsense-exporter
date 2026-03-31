package opnsense

import (
	"fmt"
	"strconv"
	"strings"
)

type KeaDhcpv4LeasesResponse struct {
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
	Rows     []struct {
		If                   string `json:"if"`
		Address              string `json:"address"`
		Hwaddr               string `json:"hwaddr"`
		ClientId             string `json:"client_id"`
		ValidLifetime        string `json:"valid_lifetime"`
		Expiration           string `json:"expire"`
		InterfaceDescription string `json:"if_descr"`
		InterfaceName        string `json:"if_name"`
		MacInfo              string `json:"mac_info"`
		IsReserved           string `json:"is_reserved"`
		Hostname             string `json:"hostname"`
		FqdnForward          string `json:"fqdn_fwd"`
		FqdnReceived         string `json:"fqdn_rev"`
		State                string `json:"state"`
		UserContext          string `json:"user_context"`
		SubnetId             string `json:"subnet_id"`
		PoolId               string `json:"pool_id"`
	}

	// This follows pattern {"name": "desc"}
	// where name is the physical interface
	// and desc is the human-readable name as set by the user
	Interfaces map[string]string
}

type KeaDhcpv4Lease struct {
	Expiration           int
	ValidLifetime        int
	Mac                  string
	ClientId             string
	Hostname             string
	Address              string
	If                   string
	InterfaceName        string
	InterfaceDescription string
}

type InterfaceInfo struct {
	Name        string
	Description string
}

type KeaDhcpv4Leases struct {
	Leases             []KeaDhcpv4Lease
	ReservedLeaseCount map[string]int
	LeaseCount         map[string]int
	Interfaces         map[string]InterfaceInfo
}

func (c *Client) FetchLeasesv4() (KeaDhcpv4Leases, *APICallError) {
	var resp KeaDhcpv4LeasesResponse
	var data KeaDhcpv4Leases

	url, ok := c.endpoints["keaDhcpv4"]
	if !ok {
		return data, &APICallError{
			Endpoint:   "keaDhcpv4",
			Message:    "endpoint not found in client endpoints",
			StatusCode: 0,
		}
	}

	err := c.do("GET", url, nil, &resp)
	if err != nil {
		return data, err
	}

	for _, row := range resp.Rows {
		fmt.Printf("Interface name: %s\n", row.InterfaceName)
		// Update total reservation count
		data.LeaseCount[row.InterfaceName] += 1

		// Update reservation count
		if strings.Compare("", row.IsReserved) != 0 {
			data.ReservedLeaseCount[row.InterfaceName] += 1
		}

		expiration, err := strconv.Atoi(row.Expiration)
		if err != nil {
			return data, &APICallError{
				Endpoint:   "keaDhcpv4",
				Message:    "expiration time is not an integer",
				StatusCode: 0,
			}
		}
		lifetime, err := strconv.Atoi(row.ValidLifetime)
		if err != nil {
			return data, &APICallError{
				Endpoint:   "keaDhcpv4",
				Message:    "valid lifetime is not an integer",
				StatusCode: 0,
			}
		}

		// Add the information in
		data.Leases = append(data.Leases, KeaDhcpv4Lease{
			If:            row.If,
			Hostname:      row.Hostname,
			Address:       row.Address,
			Mac:           row.MacInfo,
			ClientId:      row.ClientId,
			Expiration:    expiration,
			ValidLifetime: lifetime,
		})

		data.Interfaces[row.InterfaceName] = InterfaceInfo{
			Name:        row.InterfaceName,
			Description: row.InterfaceDescription,
		}
	}

	return data, nil
}
