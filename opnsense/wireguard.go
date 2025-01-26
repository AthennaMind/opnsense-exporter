package opnsense

import "log/slog"

type wireguardClientsResponse struct {
	Rows []struct {
		IfId            string  `json:"if"`
		IfType          string  `json:"type"`
		Status          string  `json:"status"`
		Name            string  `json:"name"`
		IfName          string  `json:"ifname"`
		LatestHandshake float64 `json:"latest-handshake"`
		TransferRx      float64 `json:"transfer-rx"`
		TransferTx      float64 `json:"transfer-tx"`
	} `json:"rows"`
	RowCount int `json:"rowCount"`
	Total    int `json:"total"`
	Current  int `json:"current"`
}

// WGInterfaceStatus is the custom type that represents the status of a Wireguard interface
type WGInterfaceStatus int

const (
	WGInterfaceStatusDown WGInterfaceStatus = iota
	WGInterfaceStatusUp
	WGInterfaceStatusUnknown
)

type WireguardPeers struct {
	Device          string
	DeviceName      string
	DeviceType      string
	Name            string
	LatestHandshake float64
	TransferRx      float64
	TransferTx      float64
}

type WireguardInterfaces struct {
	Device     string
	DeviceType string
	Name       string
	DeviceName string
	Status     WGInterfaceStatus
}

type WireguardClients struct {
	Peers      []WireguardPeers
	Interfaces []WireguardInterfaces
}

// parseWGInterfaceStatus parses a string status to a WGInterfaceStatus type.
func parseWGInterfaceStatus(statusTranslated string, logger *slog.Logger, originalStatus string) WGInterfaceStatus {
	switch statusTranslated {
	case "up":
		return WGInterfaceStatusUp
	case "down":
		return WGInterfaceStatusDown
	default:
		logger.Warn("unknown wireguard interface status detected", "status", originalStatus)
		return WGInterfaceStatusUnknown
	}
}

func (c *Client) FetchWireguardConfig() (WireguardClients, *APICallError) {
	var response wireguardClientsResponse
	var data WireguardClients

	url, ok := c.endpoints["wireguardClients"]
	if !ok {
		return data, &APICallError{
			Endpoint:   string(url),
			Message:    "Unable to fetch Wireguard stats",
			StatusCode: 0,
		}
	}

	if err := c.do("GET", url, nil, &response); err != nil {
		return data, err
	}

	for _, v := range response.Rows {

		if v.IfType == "interface" {
			data.Interfaces = append(data.Interfaces, WireguardInterfaces{
				Device:     v.IfId,
				DeviceType: v.IfType,
				Status:     parseWGInterfaceStatus(v.Status, c.log, v.Status),
				Name:       v.Name,
				DeviceName: v.IfName,
			})
		}
		if v.IfType == "peer" {
			data.Peers = append(data.Peers, WireguardPeers{
				DeviceType:      v.IfType,
				LatestHandshake: v.LatestHandshake,
				TransferRx:      v.TransferRx,
				TransferTx:      v.TransferTx,
				Name:            v.Name,
				DeviceName:      v.IfName,
				Device:          v.IfId,
			})
		}
	}

	return data, nil
}
