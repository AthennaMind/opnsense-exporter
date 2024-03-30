package collector

import (
	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type openVPNCollector struct {
	log       log.Logger
	subsystem string
	instance  string
	instances *prometheus.Desc
}

func init() {
	collectorInstances = append(collectorInstances, &openVPNCollector{
		subsystem: OpenVPNSubsystem,
	})
}

func (c *openVPNCollector) Name() string {
	return c.subsystem
}

func (c *openVPNCollector) Register(namespace, instanceLabel string, log log.Logger) {
	c.log = log
	c.instance = instanceLabel

	level.Debug(c.log).
		Log("msg", "Registering collector", "collector", c.Name())

	c.instances = buildPrometheusDesc(c.subsystem, "instances",
		"OpenVPN instances (1 = enabled, 0 = disabled) by role (server, client)",
		[]string{"uuid", "role", "description", "device_type"},
	)
}

func (c *openVPNCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.instances
}

func (c *openVPNCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	instances, err := client.FetchOpenVPNInstances()
	if err != nil {
		return err
	}
	for _, instance := range instances.Rows {
		ch <- prometheus.MustNewConstMetric(
			c.instances,
			prometheus.GaugeValue,
			float64(instance.Enabled),
			instance.UUID,
			instance.Role,
			instance.Description,
			instance.DevType,
			c.instance,
		)
	}
	return nil
}
