package collector

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AthennaMind/opnsense-exporter/internal/options"
	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type resourceUsageTestCollector struct {
	CollectorInstance
	client *opnsense.Client
	err    *opnsense.APICallError
}

func (c *resourceUsageTestCollector) Collect(ch chan<- prometheus.Metric) {
	c.err = c.Update(c.client, ch)
}

func TestResourceUsageCollectors(t *testing.T) {
	for _, tc := range []struct {
		name      string
		collector CollectorInstance
		path      string
		payload   string
		metrics   map[string]float64
	}{
		{
			name:      "pf states",
			collector: &pfStatesCollector{},
			path:      "/api/diagnostics/firewall/pf_states",
			payload:   `{"current":"0","limit":"1000000"}`,
			metrics:   map[string]float64{"opnsense_pf_states_current": 0, "opnsense_pf_states_limit": 1000000},
		},
		{
			name:      "mbuf",
			collector: &mbufCollector{},
			path:      "/api/diagnostics/system/system_mbuf",
			payload:   `{"mbuf-statistics":{"cluster-total":256,"cluster-max":524288,"mbuf-current":123}}`,
			metrics:   map[string]float64{"opnsense_mbuf_clusters_current": 256, "opnsense_mbuf_clusters_limit": 524288},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			responses := []struct {
				status  int
				payload string
			}{
				{http.StatusOK, tc.payload},
				{http.StatusForbidden, `{"message":"Forbidden"}`},
				{http.StatusOK, `{"result":"failed"}`},
				{http.StatusOK, tc.payload},
			}
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.path || requests >= len(responses) {
					t.Errorf("unexpected request: %s", r.URL.Path)
					http.NotFound(w, r)
					return
				}
				response := responses[requests]
				requests++
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(response.status)
				_, _ = io.WriteString(w, response.payload)
			}))
			defer server.Close()
			log := slog.New(slog.NewTextHandler(io.Discard, nil))
			client, err := opnsense.NewClient(options.OPNSenseConfig{
				Protocol: "http",
				Host:     strings.TrimPrefix(server.URL, "http://"),
			}, "test", log)
			if err != nil {
				t.Fatal(err)
			}
			tc.collector.Register(namespace, "firewall-a", log)
			adapter := &resourceUsageTestCollector{CollectorInstance: tc.collector, client: &client}
			registry := prometheus.NewPedanticRegistry()
			registry.MustRegister(adapter)
			for scrape := range responses {
				families, err := registry.Gather()
				if err != nil {
					t.Fatal(err)
				}
				if scrape == 1 || scrape == 2 {
					if adapter.err == nil || len(families) != 0 {
						t.Fatalf("failed scrape must return error without stale or zero metrics: err=%v metrics=%v", adapter.err, families)
					}
					continue
				}
				if adapter.err != nil {
					t.Fatal(adapter.err)
				}
				if len(families) != len(tc.metrics) {
					t.Fatalf("got %d metrics, want %d", len(families), len(tc.metrics))
				}
				for _, family := range families {
					want, ok := tc.metrics[family.GetName()]
					if !ok || family.GetType().String() != "GAUGE" || len(family.Metric) != 1 {
						t.Fatalf("unexpected metric family: %v", family)
					}
					metric := family.Metric[0]
					if metric.GetGauge().GetValue() != want {
						t.Errorf("%s: got %v, want %v", family.GetName(), metric.GetGauge().GetValue(), want)
					}
					if len(metric.Label) != 1 || metric.Label[0].GetName() != "opnsense_instance" || metric.Label[0].GetValue() != "firewall-a" {
						t.Errorf("unexpected labels: %v", metric.Label)
					}
				}
			}
		})
	}
}

func TestDisableResourceUsageCollectors(t *testing.T) {
	collector := &Collector{collectors: []CollectorInstance{
		&pfStatesCollector{},
		&mbufCollector{},
		&firewallCollector{subsystem: FirewallSubsystem},
	}}
	for _, disable := range []Option{WithoutPFStatesCollector(), WithoutMbufCollector()} {
		if err := disable(collector); err != nil {
			t.Fatal(err)
		}
	}
	if len(collector.collectors) != 1 || collector.collectors[0].Name() != FirewallSubsystem {
		t.Fatal("disabling resource usage must preserve existing firewall collector")
	}
}
