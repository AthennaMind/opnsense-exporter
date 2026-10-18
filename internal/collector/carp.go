package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type carpCollector struct {
	log             *slog.Logger
	vipStatus       *prometheus.Desc
	vipAdvbase      *prometheus.Desc
	vipAdvskew      *prometheus.Desc
	demotion        *prometheus.Desc
	allowed         *prometheus.Desc
	maintenanceMode *prometheus.Desc

	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &carpCollector{
		subsystem: CarpSubsystem,
	})
}

func (c *carpCollector) Name() string {
	return c.subsystem
}

func (c *carpCollector) Register(namespace, instanceLabel string, log *slog.Logger) {
	c.log = log
	c.instance = instanceLabel
	c.log.Debug("Registering collector", "collector", c.Name())

	vipLabels := []string{"interface", "vhid", "vip"}

	c.vipStatus = buildPrometheusDesc(c.subsystem, "vip_status",
		"CARP VIP status (1 = MASTER, 0 = BACKUP, 2 = INIT, 3 = DISABLED, 4 = unknown)",
		vipLabels,
	)
	c.vipAdvbase = buildPrometheusDesc(c.subsystem, "vip_advbase",
		"CARP VIP advertisement base in seconds",
		vipLabels,
	)
	c.vipAdvskew = buildPrometheusDesc(c.subsystem, "vip_advskew",
		"CARP VIP advertisement skew",
		vipLabels,
	)
	c.demotion = buildPrometheusDesc(c.subsystem, "demotion",
		"CARP demotion factor (net.inet.carp.demotion), 0 when healthy and 240 in maintenance mode",
		nil,
	)
	c.allowed = buildPrometheusDesc(c.subsystem, "allowed",
		"Whether CARP is allowed on this node (1 = allowed, 0 = temporarily disabled)",
		nil,
	)
	c.maintenanceMode = buildPrometheusDesc(c.subsystem, "maintenance_mode",
		"Whether CARP maintenance mode is enabled (1 = enabled, 0 = disabled)",
		nil,
	)
}

func (c *carpCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.vipStatus
	ch <- c.vipAdvbase
	ch <- c.vipAdvskew
	ch <- c.demotion
	ch <- c.allowed
	ch <- c.maintenanceMode
}

func (c *carpCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	data, err := client.FetchCarpStatus()
	if err != nil {
		return err
	}

	for _, vip := range data.VIPs {
		labels := []string{vip.Interface, vip.VHID, vip.VIP, c.instance}
		ch <- prometheus.MustNewConstMetric(c.vipStatus, prometheus.GaugeValue, float64(vip.Status), labels...)
		ch <- prometheus.MustNewConstMetric(c.vipAdvbase, prometheus.GaugeValue, float64(vip.Advbase), labels...)
		ch <- prometheus.MustNewConstMetric(c.vipAdvskew, prometheus.GaugeValue, float64(vip.Advskew), labels...)
	}

	if !data.HasGlobal {
		return nil
	}

	ch <- prometheus.MustNewConstMetric(c.demotion, prometheus.GaugeValue, float64(data.Demotion), c.instance)
	ch <- prometheus.MustNewConstMetric(c.allowed, prometheus.GaugeValue, float64(data.Allowed), c.instance)
	ch <- prometheus.MustNewConstMetric(c.maintenanceMode, prometheus.GaugeValue, float64(data.MaintenanceMode), c.instance)

	return nil
}
