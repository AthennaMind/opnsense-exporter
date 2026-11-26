package opnsense

import (
	"encoding/json"
	"testing"
)

func TestSystemStatusOK(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    bool
	}{
		{
			name:    "legacy format ok",
			payload: `{"System":{"status":"OK"}}`,
			want:    true,
		},
		{
			name:    "legacy format error",
			payload: `{"System":{"status":"Error"}}`,
			want:    false,
		},
		{
			name:    "25.1 metadata format ok",
			payload: `{"metadata":{"system":{"status":2,"message":"No pending messages","title":"System"},"subsystems":[]}}`,
			want:    true,
		},
		{
			name:    "25.1 metadata format error",
			payload: `{"metadata":{"system":{"status":-1,"message":"An issue was detected","title":"System"},"subsystems":[]}}`,
			want:    false,
		},
		{
			name:    "metadata format with string ok",
			payload: `{"metadata":{"system":{"status":"OK","message":"No pending messages","title":"System"},"subsystems":[]},"subsystems":[]}`,
			want:    true,
		},
		{
			name:    "metadata format with string error",
			payload: `{"metadata":{"system":{"status":"Error","message":"An issue was detected","title":"System"},"subsystems":[]},"subsystems":[]}`,
			want:    false,
		},
		{
			name:    "empty response",
			payload: `{}`,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp HealthCheckResponse
			if err := json.Unmarshal([]byte(tt.payload), &resp); err != nil {
				t.Fatalf("failed to unmarshal payload: %v", err)
			}
			if got := resp.SystemStatusOK(); got != tt.want {
				t.Errorf("SystemStatusOK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirewallStatusOK(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    bool
	}{
		{
			name:    "legacy format ok",
			payload: `{"System":{"status":"OK"},"Firewall":{"status":"OK","statusCode":2}}`,
			want:    true,
		},
		{
			name:    "legacy format error",
			payload: `{"System":{"status":"OK"},"Firewall":{"status":"Error","statusCode":0,"message":"alert"}}`,
			want:    false,
		},
		{
			name:    "25.1 metadata format ok",
			payload: `{"metadata":{"System":{"status":2},"Firewall":{"status":2}}}`,
			want:    true,
		},
		{
			name:    "25.1 metadata format error",
			payload: `{"metadata":{"System":{"status":2},"Firewall":{"status":0,"message":"alert"}}}`,
			want:    false,
		},
		{
			name:    "newer format without firewall section",
			payload: `{"metadata":{"system":{"status":2,"message":"No pending messages","title":"System"},"subsystems":[]}}`,
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp HealthCheckResponse
			if err := json.Unmarshal([]byte(tt.payload), &resp); err != nil {
				t.Fatalf("failed to unmarshal payload: %v", err)
			}
			if got := resp.FirewallStatusOK(); got != tt.want {
				t.Errorf("FirewallStatusOK() = %v, want %v", got, tt.want)
			}
		})
	}
}
