package opnsense

import (
	"fmt"
	"strconv"
)

type KeaDhcpv4LeasesRow struct {
	If                   string   `json:"if"`
	Address              string   `json:"address"`
	Hwaddr               string   `json:"hwaddr"`
	ClientId             string   `json:"client_id"`
	ValidLifetime        int      `json:"valid_lifetime"`
	Expiration           int      `json:"expire"`
	InterfaceDescription string   `json:"if_descr"`
	InterfaceName        string   `json:"if_name"`
	MacInfo              string   `json:"mac_info"`
	IsReserved           []string `json:"is_reserved"`
	Hostname             string   `json:"hostname"`
	FqdnForward          string   `json:"fqdn_fwd"`
	FqdnReceived         string   `json:"fqdn_rev"`
	State                int      `json:"state"`
	UserContext          string   `json:"user_context"`
	SubnetId             int      `json:"subnet_id"`
	PoolId               int      `json:"pool_id"`
}

type KeaDhcpv4LeasesRowAllStrings struct {
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

type KeaDhcpv4LeasesRowIntString struct {
	If                   string `json:"if"`
	Address              string `json:"address"`
	Hwaddr               string `json:"hwaddr"`
	ClientId             string `json:"client_id"`
	ValidLifetime        int    `json:"valid_lifetime"`
	Expiration           int    `json:"expire"`
	InterfaceDescription string `json:"if_descr"`
	InterfaceName        string `json:"if_name"`
	MacInfo              string `json:"mac_info"`
	IsReserved           string `json:"is_reserved"`
	Hostname             string `json:"hostname"`
	FqdnForward          string `json:"fqdn_fwd"`
	FqdnReceived         string `json:"fqdn_rev"`
	State                int    `json:"state"`
	UserContext          string `json:"user_context"`
	SubnetId             int    `json:"subnet_id"`
	PoolId               int    `json:"pool_id"`
}

type KeaDhcpv4LeasesResponseAllStrings struct {
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
	Rows     []KeaDhcpv4LeasesRowAllStrings
}

type KeaDhcpv4LeasesResponseIntString struct {
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
	Rows     []KeaDhcpv4LeasesRowIntString
}

type KeaDhcpv4LeasesResponse struct {
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
	Rows     []KeaDhcpv4LeasesRow
}

type KeaDhcpv4Lease struct {
	Expiration    int
	ValidLifetime int
	Mac           string
	MacInfo       string
	ClientId      string
	Hostname      string
	Address       string
	InterfaceName string
}

type KeaDhcpV4InterfaceInfo struct {
	Name        string
	Description string
}

type KeaDhcpv4Leases struct {
	Leases             []KeaDhcpv4Lease
	ReservedLeaseCount map[string]int
	LeaseCount         map[string]int
	Interfaces         map[string]KeaDhcpV4InterfaceInfo
}

func parseDHCPv4LeasesAllStrings(response KeaDhcpv4LeasesResponseAllStrings) (KeaDhcpv4Leases, *APICallError) {
	data := KeaDhcpv4Leases{}
	data.Interfaces = make(map[string]KeaDhcpV4InterfaceInfo)
	data.LeaseCount = make(map[string]int)
	data.ReservedLeaseCount = make(map[string]int)

	for _, row := range response.Rows {
		// Update total reservation count
		data.LeaseCount[row.InterfaceName] += 1

		// Update reservation count
		if len(row.IsReserved) > 0 {
			data.ReservedLeaseCount[row.InterfaceName] += 1
		}

		expiration, err := strconv.Atoi(row.Expiration)
		if err != nil {
			return data, &APICallError{
				Endpoint:   "keaDhcpv4",
				Message:    "failed to parse expiration",
				StatusCode: 0,
			}
		}
		lifetime, err := strconv.Atoi(row.ValidLifetime)
		if err != nil {
			return data, &APICallError{
				Endpoint:   "keaDhcpv4",
				Message:    "failed to parse valid lifetime",
				StatusCode: 0,
			}
		}
		lease := KeaDhcpv4Lease{
			InterfaceName: row.InterfaceName,
			Hostname:      row.Hostname,
			Address:       row.Address,
			Mac:           row.Hwaddr,
			ClientId:      row.ClientId,
			MacInfo:       row.MacInfo,
			Expiration:    expiration,
			ValidLifetime: lifetime,
		}

		// Add the information in
		data.Leases = append(data.Leases, lease)

		data.Interfaces[row.InterfaceName] = KeaDhcpV4InterfaceInfo{
			Name:        row.If,
			Description: row.InterfaceDescription,
		}
	}

	return data, nil
}

func parseDHCPv4LeasesStringInt(response KeaDhcpv4LeasesResponseIntString) (KeaDhcpv4Leases, *APICallError) {
	data := KeaDhcpv4Leases{}
	data.Interfaces = make(map[string]KeaDhcpV4InterfaceInfo)
	data.LeaseCount = make(map[string]int)
	data.ReservedLeaseCount = make(map[string]int)

	for _, row := range response.Rows {
		// Update total reservation count
		data.LeaseCount[row.InterfaceName] += 1

		// Update reservation count
		if len(row.IsReserved) > 0 {
			data.ReservedLeaseCount[row.InterfaceName] += 1
		}
		lease := KeaDhcpv4Lease{
			InterfaceName: row.InterfaceName,
			Hostname:      row.Hostname,
			Address:       row.Address,
			Mac:           row.Hwaddr,
			ClientId:      row.ClientId,
			MacInfo:       row.MacInfo,
			Expiration:    row.Expiration,
			ValidLifetime: row.ValidLifetime,
		}

		// Add the information in
		data.Leases = append(data.Leases, lease)

		data.Interfaces[row.InterfaceName] = KeaDhcpV4InterfaceInfo{
			Name:        row.If,
			Description: row.InterfaceDescription,
		}
	}

	return data, nil
}

func parseDHCPv4Leases(response KeaDhcpv4LeasesResponse) (KeaDhcpv4Leases, *APICallError) {
	data := KeaDhcpv4Leases{}

	data.Interfaces = make(map[string]KeaDhcpV4InterfaceInfo)
	data.LeaseCount = make(map[string]int)
	data.ReservedLeaseCount = make(map[string]int)

	for _, row := range response.Rows {
		// Update total reservation count
		data.LeaseCount[row.InterfaceName] += 1

		// Update reservation count
		if len(row.IsReserved) > 0 {
			data.ReservedLeaseCount[row.InterfaceName] += 1
		}

		expiration := row.Expiration
		lifetime := row.ValidLifetime

		// Add the information in
		data.Leases = append(data.Leases, KeaDhcpv4Lease{
			InterfaceName: row.InterfaceName,
			Hostname:      row.Hostname,
			Address:       row.Address,
			Mac:           row.Hwaddr,
			ClientId:      row.ClientId,
			Expiration:    expiration,
			ValidLifetime: lifetime,
			MacInfo:       row.MacInfo,
		})

		data.Interfaces[row.InterfaceName] = KeaDhcpV4InterfaceInfo{
			Name:        row.If,
			Description: row.InterfaceDescription,
		}
	}

	return data, nil
}

func (c *Client) FetchLeasesv4() (KeaDhcpv4Leases, *APICallError) {
	var resp KeaDhcpv4LeasesResponse
	var allStringsResp KeaDhcpv4LeasesResponseAllStrings
	var stringIntResp KeaDhcpv4LeasesResponseIntString
	var data KeaDhcpv4Leases
	apiVariant := 0

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
		apiVariant += 1
		err = c.do("GET", url, nil, &allStringsResp)
		if err != nil {
			apiVariant += 1
			err = c.do("GET", url, nil, &stringIntResp)
			if err != nil {
				return data, err
			}
		}
	}

	switch apiVariant {
	case 0:
		data, err = parseDHCPv4Leases(resp)
	case 1:
		data, err = parseDHCPv4LeasesAllStrings(allStringsResp)
	case 2:
		data, err = parseDHCPv4LeasesStringInt(stringIntResp)
	default:
		err = &APICallError{
			Endpoint:   "keaDhcpv4",
			Message:    fmt.Sprintf("unknown api variant %d", apiVariant),
			StatusCode: 0,
		}
	}
	if err != nil {
		return data, err
	}

	return data, nil
}
