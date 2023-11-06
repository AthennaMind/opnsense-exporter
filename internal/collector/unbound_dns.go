package collector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/st3ga/opnsense-exporter/opnsense"
)

type unboundDNSCollector struct {
	log       log.Logger
	subsystem string
	instance  string
	uptime    *prometheus.Desc
}

func init() {
	collectorInstances = append(collectorInstances, &unboundDNSCollector{
		subsystem: "unbound_dns",
	})
}

func (c *unboundDNSCollector) Name() string {
	return c.subsystem
}

func (c *unboundDNSCollector) Register(namespace, instanceLabel string, log log.Logger) {
	c.log = log
	c.instance = instanceLabel
	level.Debug(c.log).
		Log("msg", "Registering collector", "collector", c.Name())

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
