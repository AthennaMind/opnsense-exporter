package opnsense

import (
	"encoding/json"
	"testing"
)

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
