package opnsense

type KeaDhcpv6LeasesRow struct {
	If                    string   `json:"if"`
	Address               string   `json:"address"`
	Hwaddr                string   `json:"hwaddr"`
	Duid                  string   `json:"duid"`
	ValidLifetime         int      `json:"valid_lifetime"`
	Expiration            int      `json:"expire"`
	InterfaceDescription  string   `json:"if_descr"`
	InterfaceName         string   `json:"if_name"`
	IsReserved            []string `json:"is_reserved"`
	Hostname              string   `json:"hostname"`
	FqdnForward           string   `json:"fqdn_fwd"`
	FqdnReceived          string   `json:"fqdn_rev"`
	State                 int      `json:"state"`
	UserContext           string   `json:"user_context"`
	SubnetId              string   `json:"subnet_id"`
	PoolId                string   `json:"pool_id"`
	PreferredLifetime     int      `json:"pref_lifetime"`
	Iaid                  string   `json:"iaid"`
	PrefixLength          int      `json:"prefix_len"`
	HardwareType          string   `json:"hwtype"`
	HardwareAddressSource string   `json:"hwaddr_source"`
}

type KeaDhcpv6LeasesResponse struct {
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
	Rows     []KeaDhcpv6LeasesRow
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

func parseDHCPv6Leases(leases KeaDhcpv6LeasesResponse) (KeaDhcpv6Leases, *APICallError) {
	data := KeaDhcpv6Leases{}

	data.Interfaces = make(map[string]KeaDhcpV6InterfaceInfo)
	data.LeaseCount = make(map[string]int)
	data.ReservedLeaseCount = make(map[string]int)

	for _, row := range leases.Rows {
		// Update total reservation count
		data.LeaseCount[row.InterfaceName] += 1

		// Update reservation count
		if len(row.IsReserved) > 0 {
			data.ReservedLeaseCount[row.InterfaceName] += 1
		}

		expiration := row.Expiration
		lifetime := row.ValidLifetime
		preferredLifetime := row.PreferredLifetime
		prefixLength := row.PrefixLength

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

	data, err = parseDHCPv6Leases(resp)
	if err != nil {
		return data, err
	}

	return data, nil
}
