package opnsense

// TODO: Add client fetching

type InterfaceDetails struct {
	Device                    string `json:"device"`
	Driver                    string `json:"driver"`
	Index                     string `json:"index"`
	Flags                     string `json:"flags"`
	PromiscuousListeners      string `json:"promiscuous listeners"`
	SendQueueLength           string `json:"send queue length"`
	SendQueueMaxLength        string `json:"send queue max length"`
	SendQueueDrops            string `json:"send queue drops"`
	Type                      string `json:"type"`
	AddressLength             string `json:"address length"`
	HeaderLength              string `json:"header length"`
	LinkState                 string `json:"link state"`
	Vhid                      string `json:"vhid"`
	Datalen                   string `json:"datalen"`
	MTU                       string `json:"mtu"`
	Metric                    string `json:"metric"`
	LineRate                  string `json:"line rate"`
	PacketsReceived           string `json:"packets received"`
	PacketsTransmitted        string `json:"packets transmitted"`
	BytesReceived             string `json:"bytes received"`
	BytesTransmitted          string `json:"bytes transmitted"`
	OutputErrors              string `json:"output errors"`
	InputErrors               string `json:"input errors"`
	Collisions                string `json:"collisions"`
	MulticastsReceived        string `json:"multicasts received"`
	MulticastsTransmitted     string `json:"multicasts transmitted"`
	InputQueueDrops           string `json:"input queue drops"`
	PacketsForUnknownProtocol string `json:"packets for unknown protocol"`
	HWOffloadCapabilities     string `json:"HW offload capabilities"`
	UptimeAtAttachOrStatReset string `json:"uptime at attach or stat reset"`
	Name                      string `json:"name"`
}

// Interface is the struct returned by the OPNsense API
// when requesting the interfaces. The response is weird json
// that have the interface name as key and the interfaceDetails struct as value
type interfaceResponse struct {
	Interface map[string]InterfaceDetails `json:"interfaces"`
}

type Interface struct {
	Name                  string
	Device                string
	Type                  string
	MTU                   int
	PacketsReceived       int
	PacketsTransmitted    int
	BytesReceived         int
	BytesTransmitted      int
	MulticastsReceived    int
	MulticastsTransmitted int
	InputErrors           int
	OutputErrors          int
	Collisions            int
}

type Interfaces struct {
	Interfaces []Interface
}

func (c *Client) FetchInterfaces() (Interfaces, *APICallError) {
	var resp interfaceResponse
	var data Interfaces

	url, ok := c.endpoints["interfaces"]
	if !ok {
		return data, &APICallError{
			Endpoint:   "arp",
			Message:    "endpoint not found in client endpoints",
			StatusCode: 0,
		}
	}

	err := c.do("GET", url, nil, &resp)
	if err != nil {
		return data, err
	}

	for _, v := range resp.Interface {

		convertedValues, err := sliceIntToMapStringInt(
			[]string{
				v.MTU, v.BytesReceived, v.BytesTransmitted,
				v.PacketsReceived, v.PacketsTransmitted,
				v.MulticastsReceived, v.MulticastsTransmitted,
				v.InputErrors, v.OutputErrors,
				v.Collisions,
			},
			url,
		)

		if err != nil {
			return data, err
		}

		data.Interfaces = append(data.Interfaces, Interface{
			Name:                  v.Name,
			Device:                v.Device,
			Type:                  v.Type,
			MTU:                   convertedValues[v.MTU],
			BytesReceived:         convertedValues[v.BytesReceived],
			BytesTransmitted:      convertedValues[v.BytesTransmitted],
			PacketsReceived:       convertedValues[v.PacketsReceived],
			PacketsTransmitted:    convertedValues[v.PacketsTransmitted],
			MulticastsReceived:    convertedValues[v.MulticastsReceived],
			MulticastsTransmitted: convertedValues[v.MulticastsTransmitted],
			InputErrors:           convertedValues[v.InputErrors],
			OutputErrors:          convertedValues[v.OutputErrors],
			Collisions:            convertedValues[v.Collisions],
		})
	}

	return data, nil
}
