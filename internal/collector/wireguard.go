package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type WireguardCollector struct {
	log             *slog.Logger
	instances       *prometheus.Desc
	TransferRx      *prometheus.Desc
	TransferTx      *prometheus.Desc
	LatestHandshake *prometheus.Desc

	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &WireguardCollector{
		subsystem: WireguardSubsystem,
	})
}

func (c *WireguardCollector) Name() string {
	return c.subsystem
}

func (c *WireguardCollector) Register(namespace, instanceLabel string, log *slog.Logger) {
	c.log = log
	c.instance = instanceLabel

	c.log.Debug("Registering collector", "collector", c.Name())

	c.instances = buildPrometheusDesc(c.subsystem, "interfaces_status",
		"Wireguard interface (1 = up, 0 = down)",
		[]string{"device", "device_type", "device_name"},
	)

	c.TransferRx = buildPrometheusDesc(c.subsystem, "peer_received_bytes_total",
		"Bytes received by this wireguard peer",
		[]string{"device", "device_type", "device_name", "peer_name"},
	)

	c.TransferTx = buildPrometheusDesc(c.subsystem, "peer_transmitted_bytes_total",
		"Bytes transmitted by this wireguard peer",
		[]string{"device", "device_type", "device_name", "peer_name"},
	)

	c.LatestHandshake = buildPrometheusDesc(c.subsystem, "peer_last_handshake_seconds",
		"Last handshake by peer in seconds",
		[]string{"device", "device_type", "device_name", "peer_name"},
	)
}

func (c *WireguardCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.instances
	ch <- c.LatestHandshake
	ch <- c.TransferRx
	ch <- c.TransferTx
}

func (c *WireguardCollector) update(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(
		desc, valueType, value, labelValues...,
	)
}

func (c *WireguardCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	data, err := client.FetchWireguardConfig()
	if err != nil {
		return err
	}

	for _, instance := range data.Interfaces {
		c.update(ch, c.instances, prometheus.GaugeValue, float64(instance.Status), instance.Device, instance.DeviceType, instance.DeviceName, c.instance)
	}

	for _, instance := range data.Peers {
		c.update(ch, c.LatestHandshake, prometheus.CounterValue, float64(instance.LatestHandshake), instance.Device, instance.DeviceType, instance.DeviceName, instance.Name, c.instance)
		c.update(ch, c.TransferRx, prometheus.CounterValue, float64(instance.TransferRx), instance.Device, instance.DeviceType, instance.DeviceName, instance.Name, c.instance)
		c.update(ch, c.TransferTx, prometheus.CounterValue, float64(instance.TransferTx), instance.Device, instance.DeviceType, instance.DeviceName, instance.Name, c.instance)
	}

	return nil
}
