package opnsense

import (
	"testing"

	"github.com/prometheus/common/promslog"
)

func TestParseGatewayStatus(t *testing.T) {
	logger := promslog.NewNopLogger()
	tests := []struct {
		name string
		in   string
		want GatewayStatusType
	}{
		{"online", "Online", GatewayStatusOnline},
		{"offline", "Offline", GatewayStatusOffline},
		{"pending", "Pending", GatewayStatusPeding},
		{"packetloss", "Packetloss", GatewayStatusLoss},
		{"latency", "Latency", GatewayStatusLatency},
		{"forced", "Offline (forced)", GatewayStatusForcedDown},
		{"combined latency packetloss", "Latency, Packetloss", GatewayStatusLoss},
		{"combined online latency", "Online, Latency", GatewayStatusLatency},
		{"unknown", "Mystery", GatewayStatusUnknown},
		{"empty", "", GatewayStatusUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseGatewayStatus(tt.in, logger, tt.in)
			if got != tt.want {
				t.Errorf("parseGatewayStatus(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}
