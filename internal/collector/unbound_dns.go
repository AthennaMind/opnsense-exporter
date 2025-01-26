package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type unboundDNSCollector struct {
	log    *slog.Logger
	uptime *prometheus.Desc

	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &unboundDNSCollector{
		subsystem: UnboundDNSSubsystem,
	})
}

func (c *unboundDNSCollector) Name() string {
	return c.subsystem
}

func (c *unboundDNSCollector) Register(namespace, instanceLabel string, log *slog.Logger) {
	c.log = log
	c.instance = instanceLabel
	c.log.Debug("Registering collector", "collector", c.Name())

	c.uptime = buildPrometheusDesc(c.subsystem, "uptime_seconds",
		"Uptime of the unbound DNS service in seconds",
		nil,
	)
}

func (c *unboundDNSCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.uptime
}

func (c *unboundDNSCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	data, err := client.FetchUnboundOverview()
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		c.uptime,
		prometheus.GaugeValue,
		float64(data.UptimeSeconds),
		c.instance,
	)

	return nil
}
