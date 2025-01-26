package collector

import (
	"fmt"
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type arpTableCollector struct {
	entries   *prometheus.Desc
	log       *slog.Logger
	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &arpTableCollector{
		subsystem: ArpTableSubsystem,
	})
}

func (c *arpTableCollector) Name() string {
	return c.subsystem
}

func (c *arpTableCollector) Register(namespace, instance string, log *slog.Logger) {
	c.log = log
	c.instance = instance

	c.log.Debug("Registering collector", "collector", c.Name())

	c.entries = buildPrometheusDesc(c.subsystem, "entries",
		"Arp entries by ip, mac, hostname, interface description, type, expired and permanent",
		[]string{"ip", "mac", "hostname", "interface_description", "type", "expired", "permanent"},
	)
}

func (c *arpTableCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.entries
}

func (c *arpTableCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	data, err := client.FetchArpTable()
	if err != nil {
		return err
	}

	for _, arp := range data.Arp {
		ch <- prometheus.MustNewConstMetric(
			c.entries,
			prometheus.GaugeValue,
			1,
			arp.IP,
			arp.Mac,
			arp.Hostname,
			arp.IntfDescription,
			arp.Type,
			fmt.Sprintf("%t", arp.Expired),
			fmt.Sprintf("%t", arp.Permanent),
			c.instance,
		)
	}

	return nil
}
