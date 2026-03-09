package opnsense

import (
	"testing"

	"github.com/prometheus/common/promslog"
)

func TestParseWGPeerStatus(t *testing.T) {
	logger := promslog.NewNopLogger()

	tests := []struct {
		name     string
		status   string
		expected WGPeerStatus
	}{
		{
			name:     "Online status",
			status:   "online",
			expected: WGPeerStatusUp,
		},
		{
			name:     "Offline status",
			status:   "offline",
			expected: WGPeerStatusDown,
		},
		{
			name:     "Stale status",
			status:   "stale",
			expected: WGPeerStatusStale,
		},
		{
			name:     "Unknown status",
			status:   "something_else",
			expected: WGPeerStatusUnknown,
		},
		{
			name:     "Empty status",
			status:   "",
			expected: WGPeerStatusUnknown,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := parseWGPeerStatus(tc.status, logger, tc.status)
			if result != tc.expected {
				t.Errorf("parseWGPeerStatus(%s) = %v; want %v",
					tc.status, result, tc.expected)
			}
		})
	}
}

func TestParseWGInterfaceStatus(t *testing.T) {
	logger := promslog.NewNopLogger()

	tests := []struct {
		name     string
		status   string
		expected WGInterfaceStatus
	}{
		{
			name:     "Up status",
			status:   "up",
			expected: WGInterfaceStatusUp,
		},
		{
			name:     "Down status",
			status:   "down",
			expected: WGInterfaceStatusDown,
		},
		{
			name:     "Unknown status",
			status:   "something_else",
			expected: WGInterfaceStatusUnknown,
		},
		{
			name:     "Empty status",
			status:   "",
			expected: WGInterfaceStatusUnknown,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := parseWGInterfaceStatus(tc.status, logger, tc.status)
			if result != tc.expected {
				t.Errorf("parseWGInterfaceStatus(%s) = %v; want %v",
					tc.status, result, tc.expected)
			}
		})
	}
}

func TestWGPeerStatusValues(t *testing.T) {
	// Verify the numeric values match the expected Prometheus metric values
	tests := []struct {
		name     string
		status   WGPeerStatus
		expected int
	}{
		{
			name:     "Down equals 0",
			status:   WGPeerStatusDown,
			expected: 0,
		},
		{
			name:     "Up equals 1",
			status:   WGPeerStatusUp,
			expected: 1,
		},
		{
			name:     "Unknown equals 2",
			status:   WGPeerStatusUnknown,
			expected: 2,
		},
		{
			name:     "Stale equals 3",
			status:   WGPeerStatusStale,
			expected: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if int(tc.status) != tc.expected {
				t.Errorf("WGPeerStatus %s = %d; want %d",
					tc.name, int(tc.status), tc.expected)
			}
		})
	}
}
