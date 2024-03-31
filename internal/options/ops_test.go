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

	conf.Protocol = "https"

	if err := conf.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	conf.Protocol = "http"

	if err := conf.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
