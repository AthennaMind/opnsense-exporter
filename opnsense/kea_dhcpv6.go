package opnsense

import "strconv"

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
type KeaDhcpv6LeasesRowAllStrings struct {
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

type KeaDhcpv6LeasesResponseAllStrings struct {
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
	Rows     []KeaDhcpv6LeasesRowAllStrings
}
type KeaDhcpv6LeasesRowStringInt struct {
	If                    string `json:"if"`
	Address               string `json:"address"`
	Hwaddr                string `json:"hwaddr"`
	Duid                  string `json:"duid"`
	ValidLifetime         int    `json:"valid_lifetime"`
	Expiration            int    `json:"expire"`
	InterfaceDescription  string `json:"if_descr"`
	InterfaceName         string `json:"if_name"`
	IsReserved            string `json:"is_reserved"`
	Hostname              string `json:"hostname"`
	FqdnForward           string `json:"fqdn_fwd"`
	FqdnReceived          string `json:"fqdn_rev"`
	State                 int    `json:"state"`
	UserContext           string `json:"user_context"`
	SubnetId              string `json:"subnet_id"`
	PoolId                string `json:"pool_id"`
	PreferredLifetime     int    `json:"pref_lifetime"`
	Iaid                  string `json:"iaid"`
	PrefixLength          int    `json:"prefix_len"`
	HardwareType          string `json:"hwtype"`
	HardwareAddressSource string `json:"hwaddr_source"`
}

type KeaDhcpv6LeasesResponseStringInt struct {
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
	Rows     []KeaDhcpv6LeasesRowStringInt
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

func parseDHCPv6LeasesAllStrings(leases KeaDhcpv6LeasesResponseAllStrings) (KeaDhcpv6Leases, *APICallError) {
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

		expiration, err := strconv.Atoi(row.Expiration)
		if err != nil {
			return data, &APICallError{
				Endpoint:   "keaDhcpv6",
				Message:    "failed to parse expiration",
				StatusCode: 0,
			}
		}
		lifetime, err := strconv.Atoi(row.ValidLifetime)
		if err != nil {
			return data, &APICallError{
				Endpoint:   "keaDhcpv6",
				Message:    "failed to parse valid lifetime",
				StatusCode: 0,
			}
		}
		preferredLifetime, err := strconv.Atoi(row.PreferredLifetime)
		if err != nil {
			return data, &APICallError{
				Endpoint:   "keaDhcpv6",
				Message:    "failed to parse preferred lifetime",
				StatusCode: 0,
			}
		}
		prefixLength, err := strconv.Atoi(row.PrefixLength)
		if err != nil {
			return data, &APICallError{
				Endpoint:   "keaDhcpv6",
				Message:    "failed to parse prefix length",
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

func parseDHCPv6LeasesStringInt(leases KeaDhcpv6LeasesResponseStringInt) (KeaDhcpv6Leases, *APICallError) {
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
	var allStringsResp KeaDhcpv6LeasesResponseAllStrings
	var stringIntResp KeaDhcpv6LeasesResponseStringInt
	var data KeaDhcpv6Leases
	apiVariant := 0

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

	data, err = parseDHCPv6Leases(resp)
	if err != nil {
		return data, err
	}

	return data, nil
}
