package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type openVPNCollector struct {
	log       *slog.Logger
	instances *prometheus.Desc
	sessions  *prometheus.Desc

	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &openVPNCollector{
		subsystem: OpenVPNSubsystem,
	})
}

func (c *openVPNCollector) Name() string {
	return c.subsystem
}

func (c *openVPNCollector) Register(namespace, instanceLabel string, log *slog.Logger) {
	c.log = log
	c.instance = instanceLabel

	c.log.Debug("Registering collector", "collector", c.Name())

	c.instances = buildPrometheusDesc(c.subsystem, "instances",
		"OpenVPN instances (1 = enabled, 0 = disabled) by role (server, client)",
		[]string{"uuid", "role", "description", "device_type"},
	)
	c.sessions = buildPrometheusDesc(c.subsystem, "sessions",
		"OpenVPN session (1 = ok, 0 = not ok)",
		[]string{"description", "virtual_address", "username"},
	)
}

func (c *openVPNCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.instances
	ch <- c.sessions
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

	sessions, err := client.FetchOpenVPNSessions()
	if err != nil {
		return err
	}
	for _, session := range sessions.Rows {
		ch <- prometheus.MustNewConstMetric(
			c.sessions,
			prometheus.GaugeValue,
			float64(session.Status),
			session.Description,
			session.VirtualAddress,
			session.Username,
			c.instance,
		)
	}

	return nil
}
