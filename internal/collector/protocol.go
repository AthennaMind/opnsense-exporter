package collector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/st3ga/opnsense-exporter/opnsense"
)

type protocolCollector struct {
	log       log.Logger
	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &protocolCollector{
		subsystem: "proto_statistics",
	})
}

func (c *protocolCollector) Name() string {
	return c.subsystem
}

func (c *protocolCollector) Register(namespace, instanceLabel string, log log.Logger) {
	c.log = log
	c.instance = instanceLabel
	level.Debug(c.log).
		Log("msg", "Registering collector", "collector", c.Name())
}

func (c *protocolCollector) Describe(ch chan<- *prometheus.Desc) {

}

func (c *protocolCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {

	return nil
}
