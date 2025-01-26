package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type protocolCollector struct {
	log *slog.Logger

	tcpConnectionCountByState *prometheus.Desc
	tcpSentPackets            *prometheus.Desc
	tcpReceivedPackets        *prometheus.Desc

	arpSentRequests     *prometheus.Desc
	arpReceivedRequests *prometheus.Desc

	icmpCalls           *prometheus.Desc
	icmpSentPackets     *prometheus.Desc
	icmpDroppedByReason *prometheus.Desc

	udpDeliveredPackets  *prometheus.Desc
	udpOutputPackets     *prometheus.Desc
	udpReceivedDatagrams *prometheus.Desc
	udpDroppedByReason   *prometheus.Desc

	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &protocolCollector{
		subsystem: ProtocolSubsystem,
	})
}

func (c *protocolCollector) Name() string {
	return c.subsystem
}

func (c *protocolCollector) Register(namespace, instanceLabel string, log *slog.Logger) {
	c.log = log
	c.instance = instanceLabel
	c.log.Debug("Registering collector", "collector", c.Name())

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
	c.icmpCalls = buildPrometheusDesc(c.subsystem, "icmp_calls_total",
		"Number of ICMP calls",
		nil,
	)
	c.icmpSentPackets = buildPrometheusDesc(c.subsystem, "icmp_sent_packets_total",
		"Number of sent ICMP packets",
		nil,
	)
	c.icmpDroppedByReason = buildPrometheusDesc(c.subsystem, "icmp_dropped_by_reason_total",
		"Number of dropped ICMP packets by reason",
		[]string{"reason"},
	)
	c.udpDeliveredPackets = buildPrometheusDesc(c.subsystem, "udp_delivered_packets_total",
		"Number of delivered UDP packets",
		nil,
	)

	c.udpOutputPackets = buildPrometheusDesc(c.subsystem, "udp_output_packets_total",
		"Number of output UDP packets",
		nil,
	)

	c.udpReceivedDatagrams = buildPrometheusDesc(c.subsystem, "udp_received_datagrams_total",
		"Number of received UDP datagrams",
		nil,
	)

	c.udpDroppedByReason = buildPrometheusDesc(c.subsystem, "udp_dropped_by_reason_total",
		"Number of dropped UDP packets by reason",
		[]string{"reason"},
	)
}

func (c *protocolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.tcpConnectionCountByState
	ch <- c.tcpSentPackets
	ch <- c.tcpReceivedPackets
	ch <- c.arpSentRequests
	ch <- c.arpReceivedRequests
	ch <- c.icmpCalls
	ch <- c.icmpSentPackets
	ch <- c.icmpDroppedByReason
	ch <- c.udpDeliveredPackets
	ch <- c.udpOutputPackets
	ch <- c.udpReceivedDatagrams
	ch <- c.udpDroppedByReason
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

	ch <- prometheus.MustNewConstMetric(
		c.icmpCalls, prometheus.CounterValue, float64(data.ICMPCalls), c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.icmpSentPackets, prometheus.CounterValue, float64(data.ICMPSentPackets), c.instance,
	)
	for reason, count := range data.ICMPDroppedByReason {
		ch <- prometheus.MustNewConstMetric(
			c.icmpDroppedByReason, prometheus.GaugeValue, float64(count), reason, c.instance,
		)
	}
	ch <- prometheus.MustNewConstMetric(
		c.udpDeliveredPackets, prometheus.CounterValue, float64(data.UDPDeliveredPackets), c.instance,
	)
	ch <- prometheus.MustNewConstMetric(
		c.udpOutputPackets, prometheus.CounterValue, float64(data.UDPOutputPackets), c.instance,
	)
	ch <- prometheus.MustNewConstMetric(
		c.udpReceivedDatagrams, prometheus.CounterValue, float64(data.UDPReceivedDatagrams), c.instance,
	)
	for reason, count := range data.UDPDroppedByReason {
		ch <- prometheus.MustNewConstMetric(
			c.udpDroppedByReason, prometheus.GaugeValue, float64(count), reason, c.instance,
		)
	}
	return nil
}
