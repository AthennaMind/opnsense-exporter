package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type cronCollector struct {
	log        *slog.Logger
	jobsStatus *prometheus.Desc

	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &cronCollector{
		subsystem: CronTableSubsystem,
	})
}

func (c *cronCollector) Name() string {
	return c.subsystem
}

func (c *cronCollector) Register(namespace, instanceLabel string, log *slog.Logger) {
	c.log = log
	c.instance = instanceLabel
	c.log.Debug("Registering collector", "collector", c.Name())

	c.jobsStatus = buildPrometheusDesc(c.subsystem, "job_status",
		"Cron job status by name and description (1 = enabled, 0 = disabled)",
		[]string{"schedule", "description", "command", "origin"},
	)
}

func (c *cronCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.jobsStatus
}

func (c *cronCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	crons, err := client.FetchCronTable()
	if err != nil {
		return err
	}
	for _, cron := range crons.Cron {
		ch <- prometheus.MustNewConstMetric(
			c.jobsStatus,
			prometheus.GaugeValue,
			float64(cron.Status),
			cron.Schedule,
			cron.Description,
			cron.Command,
			cron.Origin,
			c.instance,
		)
	}
	return nil
}
