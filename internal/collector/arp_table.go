package collector

import (
	"fmt"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type arpTableCollector struct {
	log       log.Logger
	subsystem string
	instance  string
	entries   *prometheus.Desc
}

func init() {
	collectorInstances = append(collectorInstances, &arpTableCollector{
		subsystem: ArpTableSubsystem,
	})
}

func (c *arpTableCollector) Name() string {
	return c.subsystem
}

func (c *arpTableCollector) Register(namespace, instance string, log log.Logger) {
	c.log = log
	c.instance = instance

	level.Debug(c.log).
		Log("msg", "Registering collector", "collector", c.Name())

	c.entries = buildPrometheusDesc(c.subsystem, "entries",
		"Arp entries by ip, mac, hostname, interface description, type, expired and permanent",
		[]string{"ip", "mac", "hostname", "interface_description", "type", "expired", "permanent"},
	)

	// c.protocolStatistics = map[string]*prometheus.Desc{
	// 	"arpSentRequests": buildPrometheusDesc(c.subsystem, "sent_requests_total",
	// 		"Total number of sent arp requests.", nil),
	// 	"arpReceivedRequests": buildPrometheusDesc(c.subsystem, "received_requests_total",
	// 		"Total number of received arp requests", nil),
	// 	"arpSentReplies": buildPrometheusDesc(c.subsystem, "sent_replies_total",
	// 		"Total number of sent arp replies since OPNsense start.", nil),
	// 	"arpReceivedReplies": buildPrometheusDesc(c.subsystem, "received_replies_total",
	// 		"Total number of received arp replies", nil),
	// 	"arpDroppedDuplicateAddress": buildPrometheusDesc(c.subsystem, "dropped_duplicate_address_total",
	// 		"Total number of dropped arp requests due to duplicate address", nil),
	// 	"arpEntriesTimeout": buildPrometheusDesc(c.subsystem, "entries_timeout_total",
	// 		"Total number of arp entries that timed out", nil),
	// 	"arpDroppedNoEntry": buildPrometheusDesc(c.subsystem, "dropped_no_entry_total",
	// 		"Total number of dropped arp requests due to no entry", nil),
	// }

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
