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
	phase2              *prometheus.Desc
	phase2_install_time *prometheus.Desc
	phase2_bytes_in     *prometheus.Desc
	phase2_bytes_out    *prometheus.Desc
	phase2_packets_in   *prometheus.Desc
	phase2_packets_out  *prometheus.Desc
	phase2_rekey_time   *prometheus.Desc
	phase2_life_time    *prometheus.Desc

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
		[]string{"description", "name"},
	)
	c.phase1_install_time = buildPrometheusDesc(c.subsystem, "phase1_install_time",
		"IPsec phase1 install time",
		[]string{"description", "name"},
	)
	c.phase1_bytes_in = buildPrometheusDesc(c.subsystem, "phase1_bytes_in",
		"IPsec phase1 bytes in",
		[]string{"description", "name"},
	)
	c.phase1_bytes_out = buildPrometheusDesc(c.subsystem, "phase1_bytes_out",
		"IPsec phase1 bytes out",
		[]string{"description", "name"},
	)
	c.phase1_packets_in = buildPrometheusDesc(c.subsystem, "phase1_packets_in",
		"IPsec phase1 packets in",
		[]string{"description", "name"},
	)
	c.phase1_packets_out = buildPrometheusDesc(c.subsystem, "phase1_packets_out",
		"IPsec phase1 packets out",
		[]string{"description", "name"},
	)

	c.phase2_install_time = buildPrometheusDesc(c.subsystem, "phase2_install_time",
		"IPsec phase2 install time",
		[]string{"description", "name", "spi_in", "spi_out", "phase1_name"},
	)
	c.phase2_bytes_in = buildPrometheusDesc(c.subsystem, "phase2_bytes_in",
		"IPsec phase2 bytes in",
		[]string{"description", "name", "spi_in", "spi_out", "phase1_name"},
	)
	c.phase2_bytes_out = buildPrometheusDesc(c.subsystem, "phase2_bytes_out",
		"IPsec phase2 bytes out",
		[]string{"description", "name", "spi_in", "spi_out", "phase1_name"},
	)
	c.phase2_packets_in = buildPrometheusDesc(c.subsystem, "phase2_packets_in",
		"IPsec phase2 packets in",
		[]string{"description", "name", "spi_in", "spi_out", "phase1_name"},
	)
	c.phase2_packets_out = buildPrometheusDesc(c.subsystem, "phase2_packets_out",
		"IPsec phase2 packets out",
		[]string{"description", "name", "spi_in", "spi_out", "phase1_name"},
	)
	c.phase2_rekey_time = buildPrometheusDesc(c.subsystem, "phase2_rekey_time",
		"IPsec phase2 rekey time",
		[]string{"description", "name", "spi_in", "spi_out", "phase1_name"},
	)
	c.phase2_life_time = buildPrometheusDesc(c.subsystem, "phase2_life_time",
		"IPsec phase2 life time",
		[]string{"description", "name", "spi_in", "spi_out", "phase1_name"},
	)
}

func (c *ipsecCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.phase1
	ch <- c.phase1_install_time
	ch <- c.phase1_bytes_in
	ch <- c.phase1_bytes_out
	ch <- c.phase1_packets_in
	ch <- c.phase1_packets_out

	ch <- c.phase2_install_time
	ch <- c.phase2_bytes_in
	ch <- c.phase2_bytes_out
	ch <- c.phase2_packets_in
	ch <- c.phase2_packets_out
	ch <- c.phase2_rekey_time
	ch <- c.phase2_life_time
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
			phase1.Name,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.phase1_install_time,
			prometheus.GaugeValue,
			float64(phase1.InstallTime),
			phase1.Phase1desc,
			phase1.Name,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.phase1_bytes_in,
			prometheus.GaugeValue,
			float64(phase1.BytesIn),
			phase1.Phase1desc,
			phase1.Name,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.phase1_bytes_out,
			prometheus.GaugeValue,
			float64(phase1.BytesOut),
			phase1.Phase1desc,
			phase1.Name,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.phase1_packets_in,
			prometheus.GaugeValue,
			float64(phase1.PacketsIn),
			phase1.Phase1desc,
			phase1.Name,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.phase1_packets_out,
			prometheus.GaugeValue,
			float64(phase1.PacketsOut),
			phase1.Phase1desc,
			phase1.Name,
			c.instance,
		)
		for _, phase2 := range phase1.Phase2 {
			ch <- prometheus.MustNewConstMetric(
				c.phase2_install_time,
				prometheus.GaugeValue,
				float64(phase2.InstallTime),
				phase2.Phase2desc,
				phase2.Name,
				phase2.SpiIn,
				phase2.SpiOut,
				phase1.Name,
				c.instance,
			)
			ch <- prometheus.MustNewConstMetric(
				c.phase2_bytes_in,
				prometheus.GaugeValue,
				float64(phase2.BytesIn),
				phase2.Phase2desc,
				phase2.Name,
				phase2.SpiIn,
				phase2.SpiOut,
				phase1.Name,
				c.instance,
			)
			ch <- prometheus.MustNewConstMetric(
				c.phase2_bytes_out,
				prometheus.GaugeValue,
				float64(phase2.BytesOut),
				phase2.Phase2desc,
				phase2.Name,
				phase2.SpiIn,
				phase2.SpiOut,
				phase1.Name,
				c.instance,
			)
			ch <- prometheus.MustNewConstMetric(
				c.phase2_packets_in,
				prometheus.GaugeValue,
				float64(phase2.PacketsIn),
				phase2.Phase2desc,
				phase2.Name,
				phase2.SpiIn,
				phase2.SpiOut,
				phase1.Name,
				c.instance,
			)
			ch <- prometheus.MustNewConstMetric(
				c.phase2_packets_out,
				prometheus.GaugeValue,
				float64(phase2.PacketsOut),
				phase2.Phase2desc,
				phase2.Name,
				phase2.SpiIn,
				phase2.SpiOut,
				phase1.Name,
				c.instance,
			)
			ch <- prometheus.MustNewConstMetric(
				c.phase2_rekey_time,
				prometheus.GaugeValue,
				float64(phase2.RekeyTime),
				phase2.Phase2desc,
				phase2.Name,
				phase2.SpiIn,
				phase2.SpiOut,
				phase1.Name,
				c.instance,
			)
			ch <- prometheus.MustNewConstMetric(
				c.phase2_life_time,
				prometheus.GaugeValue,
				float64(phase2.LifeTime),
				phase2.Phase2desc,
				phase2.Name,
				phase2.SpiIn,
				phase2.SpiOut,
				phase1.Name,
				c.instance,
			)
		}
	}
	return nil
}
