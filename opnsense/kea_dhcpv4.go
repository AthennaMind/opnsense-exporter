package opnsense

import "fmt"

type KeaDhcpv4LeasesResponse struct {
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
	Rows     []struct {
		If                   string `json:"if"`
		Address              string `json:"address"`
		Hwaddr               string `json:"hwaddr"`
		ClientId             int    `json:"client_id"`
		ValidLifetime        int    `json:"valid_lifetime"`
		InterfaceDescription string `json:"if_descr"`
		InterfaceName        string `json:"if_name"`
		MacInfo              string `json:"mac_info"`
		IsReserved           string `json:"is_reserved"`

		// I have only seen these as "0"
		// Assuming int as default values are 0
		FqdnForward  int `json:"fqdn_fwd"`
		FqdnReceived int `json:"fqdn_rev"`

		Hostname string `json:"hostname"`

		// I have no idea what this is.
		// Default was "0" so assuming it is an integer
		State int `json:"state"`

		// I have only seen as ""
		// Out of abundance of caution leaving as string
		UserContext string `json:"user_context"`

		// This could be correlated with a little bit more work
		SubnetId int `json:"subnet_id"`

		// Suspect this takes effect if there's multiple pools in a subnet
		PoolId int `json:"pool_id"`
	}

	// This follows pattern {"name": "desc"}
	// where name is the physical interface
	// and desc is the human-readable name as set by the user
	Interfaces map[string]string
}

type KeaDhcpv4Lease struct {
}

type KeaDhcpv4Leases struct {
	Leases []KeaDhcpv4Lease
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

	fmt.Printf("%v\n", resp)

	return data, nil
}
