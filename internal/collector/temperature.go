package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type temperatureCollector struct {
	log     *slog.Logger
	celsius *prometheus.Desc

	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &temperatureCollector{
		subsystem: TemperatureSubsystem,
	})
}

func (c *temperatureCollector) Name() string {
	return c.subsystem
}

func (c *temperatureCollector) Register(namespace, instanceLabel string, log *slog.Logger) {
	c.log = log
	c.instance = instanceLabel
	c.log.Debug("Registering collector", "collector", c.Name())

	c.celsius = buildPrometheusDesc(c.subsystem, "celsius",
		"Temperature sensor reading in degrees Celsius by device and sensor type",
		[]string{"device", "device_seq", "type"},
	)
}

func (c *temperatureCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.celsius
}

func (c *temperatureCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	data, err := client.FetchTemperatures()
	if err != nil {
		return err
	}

	for _, sensor := range data.Sensors {
		ch <- prometheus.MustNewConstMetric(
			c.celsius,
			prometheus.GaugeValue,
			sensor.Celsius,
			sensor.Device,
			sensor.DeviceSeq,
			sensor.Type,
			c.instance,
		)
	}
	return nil
}
