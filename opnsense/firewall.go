package opnsense

type FirewallPFStat struct {
	InterfaceName string `json:"interface,omitempty"` // We will populate this field with the key of the map
	References    int    `json:"references"`

	In4PassPackets   int `json:"in4_pass_packets"`
	In4BlockPackets  int `json:"in4_block_packets"`
	Out4PassPackets  int `json:"out4_pass_packets"`
	Out4BlockPackets int `json:"out4_block_packets"`

	In6PassPackets   int `json:"in6_pass_packets"`
	In6BlockPackets  int `json:"in6_block_packets"`
	Out6PassPackets  int `json:"out6_pass_packets"`
	Out6BlockPackets int `json:"out6_block_packets"`
}

// firewallPFStatsResponse is the struct returned by the OPNsense API
// when requesting the firewwall statistics by interface. The response is weird json
// that have the interface name as key and the FirewallPFStats struct as value
type firewallPFStatsResponse struct {
	Interface map[string]FirewallPFStat `json:"interfaces"`
}

type FirewallPFStats struct {
	Interfaces []FirewallPFStat
}

func (c *Client) FetchPFStatsByInterface() (FirewallPFStats, *APICallError) {
	var resp firewallPFStatsResponse
	var data FirewallPFStats

	url, ok := c.endpoints["pfStatisticsByInterface"]
	if !ok {
		return data, &APICallError{
			Endpoint:   "pfStatisticsByInterface",
			Message:    "endpoint not found in client endpoints",
			StatusCode: 0,
		}
	}

	err := c.do("GET", url, nil, &resp)
	if err != nil {
		return data, err
	}

	for k, v := range resp.Interface {
		v.InterfaceName = k
		data.Interfaces = append(data.Interfaces, v)
	}
	return data, nil
}
