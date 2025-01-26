package collector

import (
	"testing"

	"github.com/AthennaMind/opnsense-exporter/internal/options"
	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/common/promslog"
)

func TestCollector(t *testing.T) {
	conf := options.OPNSenseConfig{
		Protocol: "http",
		APIKey:   "test",
	}

	client, err := opnsense.NewClient(
		conf,
		"test",
		promslog.NewNopLogger(),
	)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	collectOpts := []Option{
		WithoutArpTableCollector(),
		WithoutCronCollector(),
		WithoutUnboundCollector(),
		WithoutWireguardCollector(),
		WithoutFirewallCollector(),
	}

	collector, err := New(&client, promslog.NewNopLogger(), "test", collectOpts...)
	if err != nil {
		t.Errorf("expected no error when creating collector, got %v", err)
	}

	for _, c := range collector.collectors {
		switch c.Name() {
		case "arp_table":
			t.Errorf("expected arp_table collector to be removed")
		case "cron":
			t.Errorf("expected cron collector to be removed")
		case "unbound_dns":
			t.Errorf("expected unbound_dns collector to be removed")
		case "wireguard":
			t.Errorf("expected wireguard collector to be removed")
		case "firewall":
			t.Errorf("expected firewall collector to be removed")
		}
	}
}
