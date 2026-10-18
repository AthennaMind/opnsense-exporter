package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type mbufCollector struct {
	current  *prometheus.Desc
	limit    *prometheus.Desc
	instance string
}

func init() {
	collectorInstances = append(collectorInstances, &mbufCollector{})
}

func (c *mbufCollector) Name() string {
	return MbufSubsystem
}

func (c *mbufCollector) Register(namespace, instanceLabel string, log *slog.Logger) {
	c.instance = instanceLabel
	log.Debug("Registering collector", "collector", c.Name())
	c.current = buildPrometheusDesc(c.Name(), "clusters_current", "Current allocated mbuf clusters", nil)
	c.limit = buildPrometheusDesc(c.Name(), "clusters_limit", "Mbuf cluster allocation limit", nil)
}

func (c *mbufCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.current
	ch <- c.limit
}

func (c *mbufCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	data, err := client.FetchMbufUsage()
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(c.current, prometheus.GaugeValue, float64(data.Current), c.instance)
	ch <- prometheus.MustNewConstMetric(c.limit, prometheus.GaugeValue, float64(data.Limit), c.instance)
	return nil
}
