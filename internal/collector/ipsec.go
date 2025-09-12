package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type ipsecCollector struct {
	log    *slog.Logger
	phase1 *prometheus.Desc

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
}

func (c *ipsecCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.phase1
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
	}
	return nil
}
