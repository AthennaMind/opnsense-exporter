package options

import (
	"testing"
)

func TestOPNSenseConfig(t *testing.T) {
	conf := OPNSenseConfig{
		Protocol:  "ftp",
		Host:      "test",
		APIKey:    "test",
		APISecret: "test",
	}

	if err := conf.Validate(); err == nil {
		t.Errorf("expected invalid protocol error, got nil")
	}
}
