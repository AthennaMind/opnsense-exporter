package opnsense

import "log/slog"

type wireguardRow struct {
	IfId            string  `json:"if"`
	IfType          string  `json:"type"`
	Status          string  `json:"status"`
	Name            string  `json:"name"`
	IfName          string  `json:"ifname"`
	LatestHandshake float64 `json:"latest-handshake"`
	TransferRx      float64 `json:"transfer-rx"`
	TransferTx      float64 `json:"transfer-tx"`
	PeerStatus      string  `json:"peer-status"`
}

type wireguardClientsResponse struct {
	Rows     []wireguardRow `json:"rows"`
	RowCount int            `json:"rowCount"`
	Total    int            `json:"total"`
	Current  int            `json:"current"`
}

// WGInterfaceStatus is the custom type that represents the status of a Wireguard interface
type WGInterfaceStatus int

// WGPeerStatus is the custom type that represents the peer status
type WGPeerStatus int

const (
	WGInterfaceStatusDown WGInterfaceStatus = iota
	WGInterfaceStatusUp
	WGInterfaceStatusUnknown
)

const (
	WGPeerStatusDown WGPeerStatus = iota
	WGPeerStatusUp
	WGPeerStatusUnknown
	WGPeerStatusStale
)

type WireguardPeers struct {
	Device          string
	DeviceName      string
	DeviceType      string
	Name            string
	LatestHandshake float64
	TransferRx      float64
	TransferTx      float64
	Status          WGPeerStatus
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

// parseWGPeerStatus parses a string status to a WGPeerStatus type.
func parseWGPeerStatus(statusTranslated string, logger *slog.Logger, originalStatus string) WGPeerStatus {
	switch statusTranslated {
	case "online":
		return WGPeerStatusUp
	case "offline":
		return WGPeerStatusDown
	case "stale":
		return WGPeerStatusStale
	default:
		logger.Warn("unknown wireguard peer status detected", "status", originalStatus)
		return WGPeerStatusUnknown
	}
}

// processWireguardResponse processes wireguard API response rows and returns deduplicated data.
// Deduplication prevents Prometheus collector errors when the API returns duplicate entries.
func processWireguardResponse(rows []wireguardRow, logger *slog.Logger) WireguardClients {
	var data WireguardClients

	// Track seen entries to avoid duplicates that cause Prometheus collector errors
	seenInterfaces := make(map[string]bool)
	seenPeers := make(map[string]bool)

	for _, v := range rows {

		if v.IfType == "interface" {
			key := v.IfId + "|" + v.IfName
			if seenInterfaces[key] {
				continue
			}
			seenInterfaces[key] = true

			data.Interfaces = append(data.Interfaces, WireguardInterfaces{
				Device:     v.IfId,
				DeviceType: v.IfType,
				Status:     parseWGInterfaceStatus(v.Status, logger, v.Status),
				Name:       v.Name,
				DeviceName: v.IfName,
			})
		}
		if v.IfType == "peer" {
			key := v.IfId + "|" + v.IfName + "|" + v.Name
			if seenPeers[key] {
				continue
			}
			seenPeers[key] = true

			data.Peers = append(data.Peers, WireguardPeers{
				DeviceType:      v.IfType,
				LatestHandshake: v.LatestHandshake,
				TransferRx:      v.TransferRx,
				TransferTx:      v.TransferTx,
				Name:            v.Name,
				DeviceName:      v.IfName,
				Device:          v.IfId,
				Status:          parseWGPeerStatus(v.PeerStatus, logger, v.PeerStatus),
			})
		}
	}

	return data
}

func (c *Client) FetchWireguardConfig() (WireguardClients, *APICallError) {
	var response wireguardClientsResponse

	url, ok := c.endpoints["wireguardClients"]
	if !ok {
		return WireguardClients{}, &APICallError{
			Endpoint:   string(url),
			Message:    "Unable to fetch Wireguard stats",
			StatusCode: 0,
		}
	}

	if err := c.do("GET", url, nil, &response); err != nil {
		return WireguardClients{}, err
	}

	return processWireguardResponse(response.Rows, c.log), nil
}
