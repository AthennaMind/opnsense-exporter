package collector

import (
	"testing"

	"github.com/AthennaMind/opnsense-exporter/internal/options"
	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/go-kit/log"
)

func TestWithoutArpCollector(t *testing.T) {
	conf := options.OPNSenseConfig{
		Protocol: "http",
		APIKey:   "test",
	}

	client, err := opnsense.NewClient(
		conf,
		"test",
		log.NewNopLogger(),
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	collector, err := New(&client, log.NewNopLogger(), "test", WithoutArpTableCollector())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	for _, c := range collector.collectors {
		if c.Name() == "arp_table" {
			t.Errorf("Expected no arp collector, but it was found")
		}
	}
}
