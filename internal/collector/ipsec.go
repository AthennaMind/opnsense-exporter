package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type ipsecCollector struct {
	log                 *slog.Logger
	phase1              *prometheus.Desc
	phase1_install_time *prometheus.Desc
	phase1_bytes_in     *prometheus.Desc
	phase1_bytes_out    *prometheus.Desc
	phase1_packets_in   *prometheus.Desc
	phase1_packets_out  *prometheus.Desc

	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &ipsecCollector{
		subsystem: IPsecSubsystem,
	})
}

func (c *ipsecCollector) Name() string {
	return c.subsystem
}

func (c *ipsecCollector) Register(namespace, instanceLabel string, log *slog.Logger) {
	c.log = log
	c.instance = instanceLabel

	c.log.Debug("Registering collector", "collector", c.Name())

	c.phase1 = buildPrometheusDesc(c.subsystem, "phase1_status",
		"IPsec phase1 (1 = connected, 0 = down)",
		[]string{"description"},
	)
	c.phase1_install_time = buildPrometheusDesc(c.subsystem, "phase1_install_time",
		"IPsec phase1 install time",
		[]string{"description"},
	)
	c.phase1_bytes_in = buildPrometheusDesc(c.subsystem, "phase1_bytes_in",
		"IPsec phase1 bytes in",
		[]string{"description"},
	)
	c.phase1_bytes_out = buildPrometheusDesc(c.subsystem, "phase1_bytes_out",
		"IPsec phase1 bytes out",
		[]string{"description"},
	)
	c.phase1_packets_in = buildPrometheusDesc(c.subsystem, "phase1_packets_in",
		"IPsec phase1 packets in",
		[]string{"description"},
	)
	c.phase1_packets_out = buildPrometheusDesc(c.subsystem, "phase1_packets_out",
		"IPsec phase1 packets out",
		[]string{"description"},
	)
}

func (c *ipsecCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.phase1
	ch <- c.phase1_install_time
	ch <- c.phase1_bytes_in
	ch <- c.phase1_bytes_out
	ch <- c.phase1_packets_in
	ch <- c.phase1_packets_out
}

func (c *ipsecCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	phase1s, err := client.FetchIPsecPhase1()
	if err != nil {
		return err
	}
	for _, phase1 := range phase1s.Rows {
		ch <- prometheus.MustNewConstMetric(
			c.phase1,
			prometheus.GaugeValue,
			float64(phase1.Connected),
			phase1.Phase1desc,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.phase1_install_time,
			prometheus.GaugeValue,
			float64(phase1.InstallTime),
			phase1.Phase1desc,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.phase1_bytes_in,
			prometheus.GaugeValue,
			float64(phase1.BytesIn),
			phase1.Phase1desc,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.phase1_bytes_out,
			prometheus.GaugeValue,
			float64(phase1.BytesOut),
			phase1.Phase1desc,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.phase1_packets_in,
			prometheus.GaugeValue,
			float64(phase1.PacketsIn),
			phase1.Phase1desc,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.phase1_packets_out,
			prometheus.GaugeValue,
			float64(phase1.PacketsOut),
			phase1.Phase1desc,
			c.instance,
		)
	}
	return nil
}
