package collector

import (
	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type interfacesCollector struct {
	log                   log.Logger
	subsystem             string
	instance              string
	mtu                   *prometheus.Desc
	bytesReceived         *prometheus.Desc
	bytesTransmited       *prometheus.Desc
	multicastsTransmitted *prometheus.Desc
	multicastsReceived    *prometheus.Desc
	inputErrors           *prometheus.Desc
	outputErrors          *prometheus.Desc
	collisions            *prometheus.Desc
}

func init() {
	collectorInstances = append(collectorInstances, &interfacesCollector{
		subsystem: "interfaces",
	})
}

func (c *interfacesCollector) Name() string {
	return c.subsystem
}

func (c *interfacesCollector) Register(namespace, instanceLabel string, log log.Logger) {
	c.log = log
	c.instance = instanceLabel

	level.Debug(c.log).
		Log("msg", "Registering collector", "collector", c.Name())

	c.mtu = buildPrometheusDesc(c.subsystem, "mtu_bytes",
		"The MTU value of the interface",
		[]string{"interface", "device", "type"},
	)
	c.bytesReceived = buildPrometheusDesc(c.subsystem, "received_bytes_total",
		"Bytes received on this interface by interface name and device",
		[]string{"interface", "device", "type"},
	)
	c.bytesTransmited = buildPrometheusDesc(c.subsystem, "transmitted_bytes_total",
		"Bytes transmitted on this interface by interface name and device",
		[]string{"interface", "device", "type"},
	)
	c.multicastsReceived = buildPrometheusDesc(c.subsystem, "received_multicasts_total",
		"Multicasts received on this interface by interface name and device",
		[]string{"interface", "device", "type"},
	)
	c.multicastsTransmitted = buildPrometheusDesc(c.subsystem, "transmitted_multicasts_total",
		"Multicasts transmitted on this interface by interface name and device",
		[]string{"interface", "device", "type"},
	)
	c.inputErrors = buildPrometheusDesc(c.subsystem, "input_errors_total",
		"Input errors on this interface by interface name and device",
		[]string{"interface", "device", "type"},
	)
	c.outputErrors = buildPrometheusDesc(c.subsystem, "output_errors_total",
		"Output errors on this interface by interface name and device",
		[]string{"interface", "device", "type"},
	)
	c.collisions = buildPrometheusDesc(c.subsystem, "collisions_total",
		"Collisions on this interface by interface name and device",
		[]string{"interface", "device", "type"},
	)

}

func (c *interfacesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.mtu
	ch <- c.bytesReceived
	ch <- c.bytesTransmited
	ch <- c.multicastsReceived
	ch <- c.multicastsTransmitted
	ch <- c.inputErrors
	ch <- c.outputErrors
	ch <- c.collisions
}

func (c *interfacesCollector) update(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(
		desc, valueType, value, labelValues...,
	)
}

func (c *interfacesCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	data, err := client.FetchInterfaces()

	if err != nil {
		return err
	}

	for _, iface := range data.Interfaces {
		c.update(ch, c.mtu, prometheus.GaugeValue, float64(iface.MTU), iface.Name, iface.Device, iface.Type, c.instance)
		c.update(ch, c.bytesReceived, prometheus.CounterValue, float64(iface.BytesReceived), iface.Name, iface.Device, iface.Type, c.instance)
		c.update(ch, c.bytesTransmited, prometheus.CounterValue, float64(iface.BytesTransmitted), iface.Name, iface.Device, iface.Type, c.instance)
		c.update(ch, c.multicastsReceived, prometheus.CounterValue, float64(iface.MulticastsReceived), iface.Name, iface.Device, iface.Type, c.instance)
		c.update(ch, c.multicastsTransmitted, prometheus.CounterValue, float64(iface.MulticastsTransmitted), iface.Name, iface.Device, iface.Type, c.instance)
		c.update(ch, c.inputErrors, prometheus.CounterValue, float64(iface.InputErrors), iface.Name, iface.Device, iface.Type, c.instance)
		c.update(ch, c.outputErrors, prometheus.CounterValue, float64(iface.OutputErrors), iface.Name, iface.Device, iface.Type, c.instance)
		c.update(ch, c.collisions, prometheus.CounterValue, float64(iface.Collisions), iface.Name, iface.Device, iface.Type, c.instance)
	}

	return nil
}
