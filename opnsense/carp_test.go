package opnsense

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AthennaMind/opnsense-exporter/internal/options"
)

// Captured from OPNsense 26.7.4, extended with a BACKUP, an INIT and a
// non-CARP row.
const carpMasterPayload = `{
	"total": 4, "rowCount": 4, "current": 1,
	"carp": {"demotion": "0", "allow": "1", "maintenancemode": false, "status_msg": ""},
	"rows": [
		{"interface": "LAN", "vhid": "11", "advbase": "1", "advskew": "0", "subnet": "10.139.111.1", "status": "MASTER", "mode": "carp", "status_txt": "MASTER", "vhid_txt": "11 (freq. 1/0)"},
		{"interface": "WAN", "vhid": 8, "advbase": 1, "advskew": 100, "subnet": "192.0.2.10", "status": "BACKUP", "mode": "carp", "status_txt": "BACKUP", "vhid_txt": "8 (freq. 1/100)"},
		{"interface": "DMZ", "vhid": "12", "advbase": "1", "advskew": "0", "subnet": "2001:db8::1", "status": "INIT", "mode": "carp", "status_txt": "INIT", "vhid_txt": "12 (freq. 1/0)"},
		{"interface": "LAN", "vhid": "11", "advbase": "1", "advskew": "0", "subnet": "10.139.111.2", "status": "MASTER", "mode": "ipalias", "status_txt": "MASTER", "vhid_txt": "11"}
	]
}`

func decodeCarp(t *testing.T, payload string) Carp {
	t.Helper()
	var resp carpVIPStatusResponse
	if err := json.Unmarshal([]byte(payload), &resp); err != nil {
		t.Fatalf("failed to unmarshal carp response: %v", err)
	}
	data, err := parseCarp(resp)
	if err != nil {
		t.Fatalf("failed to parse carp response: %v", err)
	}
	return data
}

func TestParseCarp(t *testing.T) {
	data := decodeCarp(t, carpMasterPayload)

	want := []CarpVIP{
		{Interface: "LAN", VHID: "11", VIP: "10.139.111.1", Advbase: 1, Advskew: 0, Status: CarpVIPStatusMaster},
		{Interface: "WAN", VHID: "8", VIP: "192.0.2.10", Advbase: 1, Advskew: 100, Status: CarpVIPStatusBackup},
		{Interface: "DMZ", VHID: "12", VIP: "2001:db8::1", Advbase: 1, Advskew: 0, Status: CarpVIPStatusInit},
	}
	if len(data.VIPs) != len(want) {
		t.Fatalf("expected %d CARP VIPs, the ipalias row skipped, got %d: %+v", len(want), len(data.VIPs), data.VIPs)
	}
	for i, w := range want {
		if data.VIPs[i] != w {
			t.Errorf("VIP %d: got %+v, want %+v", i, data.VIPs[i], w)
		}
	}

	if !data.HasGlobal || data.Demotion != 0 || data.Allowed != 1 || data.MaintenanceMode != 0 {
		t.Errorf("unexpected global status: %+v", data)
	}
}

func TestParseCarpMaintenanceMode(t *testing.T) {
	data := decodeCarp(t, `{
		"carp": {"demotion": "240", "allow": "1", "maintenancemode": true, "status_msg": ""},
		"rows": [
			{"interface": "LAN", "vhid": "11", "advbase": "1", "advskew": "0", "subnet": "10.139.111.1", "status": "BACKUP", "mode": "carp"}
		]
	}`)

	if data.Demotion != 240 || data.MaintenanceMode != 1 || data.Allowed != 1 {
		t.Errorf("unexpected global status in maintenance mode: %+v", data)
	}
	if len(data.VIPs) != 1 || data.VIPs[0].Status != CarpVIPStatusBackup {
		t.Errorf("unexpected VIPs in maintenance mode: %+v", data.VIPs)
	}
}

func TestParseCarpDisabled(t *testing.T) {
	data := decodeCarp(t, `{
		"carp": {"demotion": "0", "allow": "0", "maintenancemode": false, "status_msg": ""},
		"rows": [
			{"interface": "LAN", "vhid": "11", "advbase": "1", "advskew": "0", "subnet": "10.139.111.1", "status": "DISABLED", "mode": "carp"}
		]
	}`)

	if data.Allowed != 0 {
		t.Errorf("expected CARP not allowed, got %+v", data)
	}
	if len(data.VIPs) != 1 || data.VIPs[0].Status != CarpVIPStatusDisabled {
		t.Errorf("unexpected VIPs with CARP disabled: %+v", data.VIPs)
	}
}

func TestParseCarpUnknownStatus(t *testing.T) {
	data := decodeCarp(t, `{
		"carp": {"demotion": "0", "allow": "1", "maintenancemode": false, "status_msg": ""},
		"rows": [
			{"interface": "LAN", "vhid": "11", "advbase": "1", "advskew": "0", "subnet": "10.139.111.1", "status": "", "mode": "carp"},
			{"interface": "WAN", "vhid": "8", "advbase": "1", "advskew": "0", "subnet": "192.0.2.10", "status": "SOMETHING_NEW", "mode": "carp"}
		]
	}`)

	if len(data.VIPs) != 2 {
		t.Fatalf("expected 2 VIPs, got %+v", data.VIPs)
	}
	for _, vip := range data.VIPs {
		if vip.Status != CarpVIPStatusUnknown {
			t.Errorf("expected unknown status for %s, got %d", vip.VIP, vip.Status)
		}
	}
	if len(data.unknownStatuses) != 2 {
		t.Errorf("expected 2 unknown statuses recorded, got %v", data.unknownStatuses)
	}
}

func TestParseCarpEmptyGlobalSection(t *testing.T) {
	data := decodeCarp(t, `{"total": 0, "rowCount": 0, "current": 1, "rows": [], "carp": []}`)

	if data.HasGlobal {
		t.Errorf("expected no global status for an empty carp section, got %+v", data)
	}
	if len(data.VIPs) != 0 {
		t.Errorf("expected no VIPs, got %+v", data.VIPs)
	}
}

func TestFetchCarpStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/diagnostics/interface/get_vip_status" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(carpMasterPayload))
	}))
	defer server.Close()

	client, err := NewClient(options.OPNSenseConfig{
		Protocol:  "http",
		Host:      strings.TrimPrefix(server.URL, "http://"),
		APIKey:    "test",
		APISecret: "test",
	}, "test", slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("failed to build client: %v", err)
	}

	data, apiErr := client.FetchCarpStatus()
	if apiErr != nil {
		t.Fatalf("unexpected error: %s", apiErr.Error())
	}
	if len(data.VIPs) != 3 || !data.HasGlobal {
		t.Errorf("unexpected carp status: %+v", data)
	}
}
