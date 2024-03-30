package collector

import (
	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type protocolCollector struct {
	log                       log.Logger
	subsystem                 string
	instance                  string
	tcpConnectionCountByState *prometheus.Desc
	tcpSentPackets            *prometheus.Desc
	tcpReceivedPackets        *prometheus.Desc
	arpSentRequests           *prometheus.Desc
	arpReceivedRequests       *prometheus.Desc
}

func init() {
	collectorInstances = append(collectorInstances, &protocolCollector{
		subsystem: ProtocolSubsystem,
	})
}

func (c *protocolCollector) Name() string {
	return c.subsystem
}

func (c *protocolCollector) Register(namespace, instanceLabel string, log log.Logger) {
	c.log = log
	c.instance = instanceLabel
	level.Debug(c.log).
		Log("msg", "Registering collector", "collector", c.Name())

	c.tcpConnectionCountByState = buildPrometheusDesc(c.subsystem, "tcp_connection_count_by_state",
		"Number of TCP connections by state",
		[]string{"state"},
	)

	c.tcpSentPackets = buildPrometheusDesc(c.subsystem, "tcp_sent_packets_total",
		"Number of sent TCP packets ",
		nil,
	)

	c.tcpReceivedPackets = buildPrometheusDesc(c.subsystem, "tcp_received_packets_total",
		"Number of received TCP packets",
		nil,
	)

	c.arpSentRequests = buildPrometheusDesc(c.subsystem, "arp_sent_requests_total",
		"Number of sent ARP requests",
		nil,
	)

	c.arpReceivedRequests = buildPrometheusDesc(c.subsystem, "arp_received_requests_total",
		"Number of received ARP requests",
		nil,
	)
}

func (c *protocolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.tcpConnectionCountByState
	ch <- c.tcpSentPackets
	ch <- c.tcpReceivedPackets
	ch <- c.arpSentRequests
	ch <- c.arpReceivedRequests
}

func (c *protocolCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	data, err := client.FetchProtocolStatistics()
	if err != nil {
		return err
	}

	for state, count := range data.TCPConnectionCountByState {
		ch <- prometheus.MustNewConstMetric(
			c.tcpConnectionCountByState, prometheus.GaugeValue, float64(count), state, c.instance,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		c.tcpSentPackets, prometheus.CounterValue, float64(data.TCPSentPackets), c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.tcpReceivedPackets, prometheus.CounterValue, float64(data.TCPReceivedPackets), c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.arpSentRequests, prometheus.CounterValue, float64(data.ARPSentRequests), c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.arpReceivedRequests, prometheus.CounterValue, float64(data.ARPReceivedRequests), c.instance,
	)

	return nil
}
