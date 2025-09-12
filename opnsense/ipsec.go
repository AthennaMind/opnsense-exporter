package opnsense

type ipsecSearchResponse struct {
	Rows []struct {
		Phase1desc string `json:"phase1desc"`
		Connected  bool   `json:"connected"`
	} `json:"rows"`
	RowCount int `json:"rowCount"`
	Total    int `json:"total"`
	Current  int `json:"current"`
}

type IPsec struct {
	Phase1desc string
	Connected  int
}

type IPsecPhase1 struct {
	Rows []IPsec
}

func (c *Client) FetchIPsecPhase1() (IPsecPhase1, *APICallError) {
	var resp ipsecSearchResponse
	var data IPsecPhase1

	url, ok := c.endpoints["ipsecPhase1"]
	if !ok {
		return data, &APICallError{
			Endpoint:   "ipsecPhase1",
			Message:    "endpoint not found in client endpoints",
			StatusCode: 0,
		}
	}

	if err := c.do("GET", url, nil, &resp); err != nil {
		return data, err
	}

	for _, v := range resp.Rows {
		data.Rows = append(data.Rows, IPsec{
			Phase1desc: v.Phase1desc,
			Connected:  parseBoolToInt(v.Connected),
		})
	}

	return data, nil
}
