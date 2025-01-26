package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type firewallCollector struct {
	log                 *slog.Logger
	inIPv4PassPackets   *prometheus.Desc
	outIPv4PassPackets  *prometheus.Desc
	inIPv4BlockPackets  *prometheus.Desc
	outIPv4BlockPackets *prometheus.Desc

	inIPv6PassPackets   *prometheus.Desc
	outIPv6PassPackets  *prometheus.Desc
	inIPv6BlockPackets  *prometheus.Desc
	outIPv6BlockPackets *prometheus.Desc

	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &firewallCollector{
		subsystem: FirewallSubsystem,
	})
}

func (c *firewallCollector) Name() string {
	return c.subsystem
}

func (c *firewallCollector) Register(namespace, instanceLabel string, log *slog.Logger) {
	c.log = log
	c.instance = instanceLabel
	c.log.Debug("Registering collector", "collector", c.Name())

	c.inIPv4PassPackets = buildPrometheusDesc(c.subsystem, "in_ipv4_pass_packets",
		"The number of IPv4 incoming packets that were allowed to pass through the firewall by interface",
		[]string{"interface"},
	)

	c.outIPv4PassPackets = buildPrometheusDesc(c.subsystem, "out_ipv4_pass_packets",
		"The number of IPv4 outgoing packets that were allowed to pass through the firewall by interface",
		[]string{"interface"},
	)

	c.inIPv4BlockPackets = buildPrometheusDesc(c.subsystem, "in_ipv4_block_packets",
		"The number of IPv4 incoming packets that were blocked by the firewall by interface",
		[]string{"interface"},
	)

	c.outIPv4BlockPackets = buildPrometheusDesc(c.subsystem, "out_ipv4_block_packets",
		"The number of IPv4 outgoing packets that were blocked by the firewall by interface",
		[]string{"interface"},
	)

	c.inIPv6PassPackets = buildPrometheusDesc(c.subsystem, "in_ipv6_pass_packets",
		"The number of IPv6 incoming packets that were allowed to pass through the firewall by interface",
		[]string{"interface"},
	)

	c.outIPv6PassPackets = buildPrometheusDesc(c.subsystem, "out_ipv6_pass_packets",
		"The number of IPv6 outgoing packets that were allowed to pass through the firewall by interface",
		[]string{"interface"},
	)

	c.inIPv6BlockPackets = buildPrometheusDesc(c.subsystem, "in_ipv6_block_packets",
		"The number of IPv6 incoming packets that were blocked by the firewall by interface",
		[]string{"interface"},
	)

	c.outIPv6BlockPackets = buildPrometheusDesc(c.subsystem, "out_ipv6_block_packets",
		"The number of IPv6 outgoing packets that were blocked by the firewall by interface",
		[]string{"interface"},
	)
}

func (c *firewallCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.inIPv4PassPackets
	ch <- c.outIPv4PassPackets
	ch <- c.inIPv4BlockPackets
	ch <- c.outIPv4BlockPackets

	ch <- c.inIPv6PassPackets
	ch <- c.outIPv6PassPackets
	ch <- c.inIPv6BlockPackets
	ch <- c.outIPv6BlockPackets
}

func (c *firewallCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	data, err := client.FetchPFStatsByInterface()
	if err != nil {
		return err
	}

	for _, v := range data.Interfaces {
		metricsValueMapping := map[*prometheus.Desc]int{
			c.inIPv4PassPackets:   v.In4PassPackets,
			c.outIPv4PassPackets:  v.Out4PassPackets,
			c.inIPv4BlockPackets:  v.In4BlockPackets,
			c.outIPv4BlockPackets: v.Out4BlockPackets,
			c.inIPv6PassPackets:   v.In6PassPackets,
			c.outIPv6PassPackets:  v.Out6PassPackets,
			c.inIPv6BlockPackets:  v.In6BlockPackets,
			c.outIPv6BlockPackets: v.Out6BlockPackets,
		}
		for metric, value := range metricsValueMapping {
			ch <- prometheus.MustNewConstMetric(
				metric,
				prometheus.GaugeValue,
				float64(value),
				v.InterfaceName,
				c.instance,
			)
		}
	}
	return nil
}
