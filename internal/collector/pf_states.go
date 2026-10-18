package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type pfStatesCollector struct {
	current  *prometheus.Desc
	limit    *prometheus.Desc
	instance string
}

func init() {
	collectorInstances = append(collectorInstances, &pfStatesCollector{})
}

func (c *pfStatesCollector) Name() string {
	return PFStatesSubsystem
}

func (c *pfStatesCollector) Register(namespace, instanceLabel string, log *slog.Logger) {
	c.instance = instanceLabel
	log.Debug("Registering collector", "collector", c.Name())
	c.current = buildPrometheusDesc(c.Name(), "current", "Current PF state table entries", nil)
	c.limit = buildPrometheusDesc(c.Name(), "limit", "PF state table entry limit", nil)
}

func (c *pfStatesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.current
	ch <- c.limit
}

func (c *pfStatesCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	data, err := client.FetchPFStates()
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(c.current, prometheus.GaugeValue, float64(data.Current), c.instance)
	ch <- prometheus.MustNewConstMetric(c.limit, prometheus.GaugeValue, float64(data.Limit), c.instance)
	return nil
}
