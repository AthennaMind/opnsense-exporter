package opnsense

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

const fetchCarpPayload = `{"current":1,"rowCount":-1,"sort":{},"searchPhrase":""}`

type carpVIPStatusResponse struct {
	Rows []struct {
		Interface string  `json:"interface"`
		VHID      flexInt `json:"vhid"`
		Advbase   flexInt `json:"advbase"`
		Advskew   flexInt `json:"advskew"`
		Subnet    string  `json:"subnet"`
		Status    string  `json:"status"`
		Mode      string  `json:"mode"`
	} `json:"rows"`
	// Upstream falls back to json_encode([]) when configd fails,
	// which yields an empty array instead of an object.
	Carp json.RawMessage `json:"carp"`
}

type carpGlobalStatus struct {
	Demotion        flexInt `json:"demotion"`
	Allow           flexInt `json:"allow"`
	MaintenanceMode bool    `json:"maintenancemode"`
}

type CarpVIPStatus int

const (
	CarpVIPStatusBackup CarpVIPStatus = iota
	CarpVIPStatusMaster
	CarpVIPStatusInit
	CarpVIPStatusDisabled
	CarpVIPStatusUnknown
)

func parseCarpVIPStatus(status string) (CarpVIPStatus, bool) {
	switch strings.ToUpper(status) {
	case "MASTER":
		return CarpVIPStatusMaster, true
	case "BACKUP":
		return CarpVIPStatusBackup, true
	case "INIT":
		return CarpVIPStatusInit, true
	case "DISABLED":
		return CarpVIPStatusDisabled, true
	default:
		return CarpVIPStatusUnknown, false
	}
}

type CarpVIP struct {
	Interface string
	VHID      string
	VIP       string
	Advbase   int
	Advskew   int
	Status    CarpVIPStatus
}

type Carp struct {
	VIPs []CarpVIP
	// HasGlobal is false when the response carried no usable carp section.
	HasGlobal       bool
	Demotion        int
	Allowed         int
	MaintenanceMode int

	unknownStatuses []string
}

func parseCarp(resp carpVIPStatusResponse) (Carp, error) {
	var data Carp

	for _, row := range resp.Rows {
		if row.Mode != "carp" {
			continue
		}
		status, ok := parseCarpVIPStatus(row.Status)
		if !ok {
			data.unknownStatuses = append(data.unknownStatuses, row.Status)
		}
		data.VIPs = append(data.VIPs, CarpVIP{
			Interface: row.Interface,
			VHID:      strconv.Itoa(int(row.VHID)),
			VIP:       row.Subnet,
			Advbase:   int(row.Advbase),
			Advskew:   int(row.Advskew),
			Status:    status,
		})
	}

	raw := bytes.TrimSpace(resp.Carp)
	if len(raw) > 0 && raw[0] == '{' {
		var global carpGlobalStatus
		if err := json.Unmarshal(raw, &global); err != nil {
			return data, err
		}
		data.HasGlobal = true
		data.Demotion = int(global.Demotion)
		data.Allowed = parseBoolToInt(global.Allow != 0)
		data.MaintenanceMode = parseBoolToInt(global.MaintenanceMode)
	}

	return data, nil
}

func (c *Client) FetchCarpStatus() (Carp, *APICallError) {
	var resp carpVIPStatusResponse
	var data Carp

	path, ok := c.endpoints["carpVIPStatus"]
	if !ok {
		return data, &APICallError{
			Endpoint:   "carpVIPStatus",
			Message:    "endpoint not found",
			StatusCode: 0,
		}
	}

	if err := c.do("POST", path, strings.NewReader(fetchCarpPayload), &resp); err != nil {
		return data, err
	}

	data, err := parseCarp(resp)
	if err != nil {
		return data, &APICallError{
			Endpoint:   "carpVIPStatus",
			Message:    "failed to decode carp section: " + err.Error(),
			StatusCode: 0,
		}
	}

	if len(data.unknownStatuses) > 0 {
		c.log.Debug("carp vips with an unrecognized status reported as unknown",
			"component", "opnsense-client",
			"statuses", data.unknownStatuses,
		)
	}

	return data, nil
}
