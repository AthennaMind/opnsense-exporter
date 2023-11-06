package collector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/st3ga/opnsense-exporter/opnsense"
)

type servicesCollector struct {
	log             log.Logger
	subsystem       string
	instance        string
	services        *prometheus.Desc
	servicesRunning *prometheus.Desc
	servicesStopped *prometheus.Desc
}

func init() {
	collectorInstances = append(collectorInstances, &servicesCollector{
		subsystem: "services",
	})
}

func (c *servicesCollector) Name() string {
	return c.subsystem
}

func (c *servicesCollector) Register(namespace, instanceLabel string, log log.Logger) {
	c.log = log
	c.instance = instanceLabel
	level.Debug(c.log).
		Log("msg", "Registering collector", "collector", c.Name())

	c.services = buildPrometheusDesc(c.subsystem, "status",
		"Service status by name and description (1 = running, 0 = stopped)",
		[]string{"name", "description"},
	)

	c.servicesRunning = buildPrometheusDesc(c.subsystem, "running_total",
		"Total number of running services",
		nil,
	)

	c.servicesStopped = buildPrometheusDesc(c.subsystem, "stopped_total",
		"Total number of stopped services",
		nil,
	)

}

func (c *servicesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.services
	ch <- c.servicesRunning
	ch <- c.servicesStopped
}

func (c *servicesCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {

	services, err := client.FetchServices()

	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.servicesRunning, prometheus.GaugeValue,
		float64(services.TotalRunning),
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.servicesStopped, prometheus.GaugeValue,
		float64(services.TotalStopped),
		c.instance,
	)

	for _, service := range services.Services {
		ch <- prometheus.MustNewConstMetric(
			c.services, prometheus.GaugeValue,
			float64(service.Status),
			service.Name,
			service.Description,
			c.instance,
		)

	}

	return nil
}
