package opnsense

import (
	"encoding/json"
	"testing"
)

func TestFlexInt(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    flexInt
		wantErr bool
	}{
		{name: "bare number", in: `86400`, want: 86400},
		{name: "quoted number", in: `"86400"`, want: 86400},
		{name: "empty string", in: `""`, want: 0},
		{name: "null", in: `null`, want: 0},
		{name: "zero string", in: `"0"`, want: 0},
		{name: "infinite lease", in: `"4294967295"`, want: 4294967295},
		{name: "non-numeric string", in: `"abc"`, wantErr: true},
		{name: "array", in: `["x"]`, wantErr: true},
		{name: "overflow", in: `"99999999999999999999"`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f flexInt
			err := json.Unmarshal([]byte(tt.in), &f)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q, got value %d", tt.in, f)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.in, err)
			}
			if f != tt.want {
				t.Errorf("flexInt(%q) = %d, want %d", tt.in, f, tt.want)
			}
		})
	}
}

func TestFlexStringSlice(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantLen int
	}{
		{name: "bare string reserved", in: `"hwaddr"`, wantLen: 1},
		{name: "empty string not reserved", in: `""`, wantLen: 0},
		{name: "null", in: `null`, wantLen: 0},
		{name: "array one", in: `["hwaddr"]`, wantLen: 1},
		{name: "array empty", in: `[]`, wantLen: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s flexStringSlice
			if err := json.Unmarshal([]byte(tt.in), &s); err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.in, err)
			}
			if len(s) != tt.wantLen {
				t.Errorf("flexStringSlice(%q) len = %d, want %d", tt.in, len(s), tt.wantLen)
			}
		})
	}
}

// TestKeaDhcpv4RowAllStringsJSON decodes a lease row in the all-strings shape
// that OPNsense 26.1.x returns from api/kea/leases4/search -- every numeric
// field quoted and is_reserved as a bare string. This is the exact payload
// that made the upstream int/[]string structs fail with
// "cannot unmarshal string into Go struct field ... valid_lifetime of type int".
func TestKeaDhcpv4RowAllStringsJSON(t *testing.T) {
	const payload = `{
		"if": "igc0", "address": "10.0.0.1", "hwaddr": "c8:9e:43:4e:9f:f6",
		"client_id": "", "valid_lifetime": "86400", "expire": "1781118374",
		"subnet_id": "1", "fqdn_fwd": "0", "fqdn_rev": "0", "hostname": "ap",
		"state": "0", "user_context": "", "pool_id": "0", "if_descr": "LAN",
		"if_name": "lan", "mac_info": "NETGEAR", "is_reserved": "hwaddr"
	}`

	var row KeaDhcpv4LeasesRow
	if err := json.Unmarshal([]byte(payload), &row); err != nil {
		t.Fatalf("failed to unmarshal all-strings lease row: %v", err)
	}
	if row.ValidLifetime != 86400 {
		t.Errorf("ValidLifetime = %d, want 86400", row.ValidLifetime)
	}
	if row.Expiration != 1781118374 {
		t.Errorf("Expiration = %d, want 1781118374", row.Expiration)
	}
	if row.State != 0 || row.SubnetId != 1 || row.PoolId != 0 {
		t.Errorf("state/subnet/pool = %d/%d/%d, want 0/1/0", row.State, row.SubnetId, row.PoolId)
	}
	if len(row.IsReserved) != 1 {
		t.Errorf("IsReserved len = %d, want 1 (reserved)", len(row.IsReserved))
	}
}
