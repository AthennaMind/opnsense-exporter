package opnsense

import "strconv"

type ipsecSearchResponse struct {
	Rows []struct {
		Phase1desc  string `json:"phase1desc"`
		Connected   bool   `json:"connected"`
		IkeId       string `json:"ikeid"`
		Name        string `json:"name"`
		InstallTime string `json:"install-time"`
		BytesIn     int    `json:"bytes-in"`
		BytesOut    int    `json:"bytes-out"`
		PacketsIn   int    `json:"packets-in"`
		PacketsOut  int    `json:"packets-out"`
	} `json:"rows"`
	RowCount int `json:"rowCount"`
	Total    int `json:"total"`
	Current  int `json:"current"`
}

type IPsec struct {
	Phase1desc  string
	Connected   int
	IkeId       string
	Name        string
	InstallTime int
	BytesIn     int
	BytesOut    int
	PacketsIn   int
	PacketsOut  int
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

	installTime, err := strconv.Atoi(resp.Rows[0].InstallTime)
	if err != nil {
		installTime = 0
	}

	for _, v := range resp.Rows {
		data.Rows = append(data.Rows, IPsec{
			Phase1desc:  v.Phase1desc,
			IkeId:       v.IkeId,
			Name:        v.Name,
			InstallTime: installTime,
			BytesIn:     v.BytesIn,
			BytesOut:    v.BytesOut,
			PacketsIn:   v.PacketsIn,
			PacketsOut:  v.PacketsOut,
			Connected:   parseBoolToInt(v.Connected),
		})
	}

	return data, nil
}
