package opnsense

type systemInfoResponse struct {
	System  string   `json:"system"`
	Plugins []string `json:"plugins"`
	Data    struct {
		Interfaces []struct {
			Inpkts       string `json:"inpkts"`
			Outpkts      string `json:"outpkts"`
			Inbytes      string `json:"inbytes"`
			Outbytes     string `json:"outbytes"`
			InbytesFrmt  string `json:"inbytes_frmt"`
			OutbytesFrmt string `json:"outbytes_frmt"`
			Inerrs       string `json:"inerrs"`
			Outerrs      string `json:"outerrs"`
			Collisions   string `json:"collisions"`
			Descr        string `json:"descr"`
			Name         string `json:"name"`
			Status       string `json:"status"`
			Ipaddr       string `json:"ipaddr"`
			Media        string `json:"media"`
		} `json:"interfaces"`
		System struct {
			Versions []string `json:"versions"`
			CPU      struct {
				Used          string   `json:"used"`
				User          string   `json:"user"`
				Nice          string   `json:"nice"`
				Sys           string   `json:"sys"`
				Intr          string   `json:"intr"`
				Idle          string   `json:"idle"`
				Model         string   `json:"model"`
				Cpus          string   `json:"cpus"`
				Cores         string   `json:"cores"`
				MaxFreq       string   `json:"max.freq"`
				CurFreq       string   `json:"cur.freq"`
				FreqTranslate string   `json:"freq_translate"`
				Load          []string `json:"load"`
			} `json:"cpu"`
			DateFrmt string `json:"date_frmt"`
			DateTime string `json:"date_time"`
			Uptime   string `json:"uptime"`
			Config   struct {
				LastChange     string `json:"last_change"`
				LastChangeFrmt string `json:"last_change_frmt"`
			} `json:"config"`
			Kernel struct {
				Pf struct {
					Maxstates string `json:"maxstates"`
					States    string `json:"states"`
				} `json:"pf"`
				Mbuf struct {
					Total string `json:"total"`
					Max   string `json:"max"`
				} `json:"mbuf"`
				Memory struct {
					Total  string `json:"total"`
					Used   string `json:"used"`
					Arc    string `json:"arc"`
					ArcTxt string `json:"arc_txt"`
				} `json:"memory"`
			} `json:"kernel"`
			Disk struct {
				Swap []struct {
					Device string `json:"device"`
					Total  string `json:"total"`
					Used   string `json:"used"`
				} `json:"swap"`
				Devices []struct {
					Device     string `json:"device"`
					Type       string `json:"type"`
					Size       string `json:"size"`
					Used       string `json:"used"`
					Available  string `json:"available"`
					Capacity   string `json:"capacity"`
					Mountpoint string `json:"mountpoint"`
				} `json:"devices"`
			} `json:"disk"`
			Firmware string `json:"firmware"`
		} `json:"system"`
		Temperature []struct {
			Device         string `json:"device"`
			DeviceSeq      string `json:"device_seq"`
			Temperature    string `json:"temperature"`
			Type           string `json:"type"`
			TypeTranslated string `json:"type_translated"`
		} `json:"temperature"`
	} `json:"data"`
}

type Temperature struct {
	Device                string
	DeviceSeq             string
	Type                  string
	TemperatureCelsuis    int64
	TemperatureFahrenheit float32
}

type SystemInfo struct {
	Temperature []Temperature
}

func (c *Client) FetchSystemInfo() (SystemInfo, *APICallError) {
	var resp systemInfoResponse
	var data SystemInfo

	url, ok := c.endpoints["systemInfo"]

	if !ok {
		return data, &APICallError{
			Endpoint:   "systemInfo",
			Message:    "endpoint not found",
			StatusCode: 0,
		}
	}

	if err := c.do("GET", url, nil, &resp); err != nil {
		return data, err
	}

	for _, v := range resp.Data.Temperature {
		celsius, err := parseStringToInt(v.Temperature, url)
		if err != nil {
			return data, err
		}
		data.Temperature = append(data.Temperature, Temperature{
			Device:                v.Device,
			DeviceSeq:             v.DeviceSeq,
			Type:                  v.Type,
			TemperatureCelsuis:    celsius,
			TemperatureFahrenheit: (float32(celsius) * 1.8) + 32,
		})
	}

	return data, nil
}
