package collector

import (
	"testing"

	"github.com/go-kit/log"
	"github.com/st3ga/opnsense-exporter/opnsense"
)

func TestWithoutArpCollector(t *testing.T) {
	client, err := opnsense.NewClient(
		"test",
		"test",
		"test",
		"test",
		"test",
		false,
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
