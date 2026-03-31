package opnsense

import (
	"strconv"
	"strings"
)

type KeaDhcpv6LeasesResponse struct {
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
	Rows     []struct {
		If                    string `json:"if"`
		Address               string `json:"address"`
		Hwaddr                string `json:"hwaddr"`
		Duid                  string `json:"duid"`
		ValidLifetime         string `json:"valid_lifetime"`
		Expiration            string `json:"expire"`
		InterfaceDescription  string `json:"if_descr"`
		InterfaceName         string `json:"if_name"`
		IsReserved            string `json:"is_reserved"`
		Hostname              string `json:"hostname"`
		FqdnForward           string `json:"fqdn_fwd"`
		FqdnReceived          string `json:"fqdn_rev"`
		State                 string `json:"state"`
		UserContext           string `json:"user_context"`
		SubnetId              string `json:"subnet_id"`
		PoolId                string `json:"pool_id"`
		PreferredLifetime     string `json:"pref_lifetime"`
		Iaid                  string `json:"iaid"`
		PrefixLength          string `json:"prefix_len"`
		HardwareType          string `json:"hwtype"`
		HardwareAddressSource string `json:"hwaddr_source"`
	}

	// This follows pattern {"name": "desc"}
	// where name is the physical interface
	// and desc is the human-readable name as set by the user
	Interfaces map[string]string
}

type KeaDhcpv6Lease struct {
	Expiration           int
	ValidLifetime        int
	PreferredLifetime    int
	Hwaddr               string
	Duid                 string
	Hostname             string
	Address              string
	PrefixLength         int
	If                   string
	InterfaceName        string
	InterfaceDescription string
}

type KeaDhcpV6InterfaceInfo struct {
	Name        string
	Description string
}

type KeaDhcpv6Leases struct {
	Leases             []KeaDhcpv6Lease
	ReservedLeaseCount map[string]int
	LeaseCount         map[string]int
	Interfaces         map[string]KeaDhcpV6InterfaceInfo
}

func (c *Client) FetchLeasesv6() (KeaDhcpv6Leases, *APICallError) {
	var resp KeaDhcpv6LeasesResponse
	var data KeaDhcpv6Leases

	url, ok := c.endpoints["keaDhcpv6"]
	if !ok {
		return data, &APICallError{
			Endpoint:   "keaDhcpv6",
			Message:    "endpoint not found in client endpoints",
			StatusCode: 0,
		}
	}

	err := c.do("GET", url, nil, &resp)
	if err != nil {
		return data, err
	}

	data.Interfaces = make(map[string]KeaDhcpV6InterfaceInfo)
	data.LeaseCount = make(map[string]int)
	data.ReservedLeaseCount = make(map[string]int)

	for _, row := range resp.Rows {
		// Update total reservation count
		data.LeaseCount[row.InterfaceName] += 1

		// Update reservation count
		if strings.Compare("", row.IsReserved) != 0 {
			data.ReservedLeaseCount[row.InterfaceName] += 1
		}

		expiration, err := strconv.Atoi(row.Expiration)
		if err != nil {
			return data, &APICallError{
				Endpoint:   "keaDhcpv6",
				Message:    "expiration time is not an integer",
				StatusCode: 0,
			}
		}
		lifetime, err := strconv.Atoi(row.ValidLifetime)
		if err != nil {
			return data, &APICallError{
				Endpoint:   "keaDhcpv6",
				Message:    "valid lifetime is not an integer",
				StatusCode: 0,
			}
		}
		preferredLifetime, err := strconv.Atoi(row.PreferredLifetime)
		if err != nil {
			return data, &APICallError{
				Endpoint:   "keaDhcpv6",
				Message:    "preferred lifetime is not an integer",
				StatusCode: 0,
			}
		}
		prefixLength, err := strconv.Atoi(row.PrefixLength)
		if err != nil {
			return data, &APICallError{
				Endpoint:   "keaDhcpv6",
				Message:    "prefix length is not an integer",
				StatusCode: 0,
			}
		}

		// Add the information in
		data.Leases = append(data.Leases, KeaDhcpv6Lease{
			InterfaceName:     row.InterfaceName,
			Hostname:          row.Hostname,
			Address:           row.Address,
			PrefixLength:      prefixLength,
			Hwaddr:            row.Hwaddr,
			Duid:              row.Duid,
			Expiration:        expiration,
			PreferredLifetime: preferredLifetime,
			ValidLifetime:     lifetime,
		})

		data.Interfaces[row.InterfaceName] = KeaDhcpV6InterfaceInfo{
			Name:        row.If,
			Description: row.InterfaceDescription,
		}
	}

	return data, nil
}
