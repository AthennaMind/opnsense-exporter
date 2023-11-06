package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/st3ga/opnsense-exporter/opnsense"
)

type gatewaysCollector struct {
	log            log.Logger
	subsystem      string
	instance       string
	status         *prometheus.Desc
	lossPercentage *prometheus.Desc
	rtt            *prometheus.Desc
	rttd           *prometheus.Desc
}

func init() {
	collectorInstances = append(collectorInstances, &gatewaysCollector{
		subsystem: "gateways",
	})
}

func (c *gatewaysCollector) Name() string {
	return c.subsystem
}

func (c *gatewaysCollector) Register(namespace, instanceLabel string, log log.Logger) {
	c.log = log
	c.instance = instanceLabel
	c.status = buildPrometheusDesc(c.subsystem, "status",
		"Status of the gateway by name and address (1 = up, 0 = down, 2 = unkown)",
		[]string{"name", "address"},
	)
	c.lossPercentage = buildPrometheusDesc(
		c.subsystem, "loss_percentage",
		"The current gateway loss percentage by name and address",
		[]string{"name", "adress"},
	)
	c.rtt = buildPrometheusDesc(
		c.subsystem, "rtt_milliseconds",
		"RTT is the average (mean) of the round trip time in milliseconds by name and address",
		[]string{"name", "adress"},
	)
	c.rttd = buildPrometheusDesc(
		c.subsystem, "rttd_milliseconds",
		"RTTd is the standard deviation of the round trip time in milliseconds by name and address",
		[]string{"name", "adress"},
	)
}

func (c *gatewaysCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.status
	ch <- c.lossPercentage
	ch <- c.rtt
}

func (c *gatewaysCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	data, err := client.FetchGateways()
	if err != nil {
		return err
	}
	for _, v := range data.Gateways {
		ch <- prometheus.MustNewConstMetric(
			c.status,
			prometheus.GaugeValue,
			float64(v.Status),
			v.Name,
			v.Address,
			c.instance,
		)
		if v.LossPercentage != -1 {
			ch <- prometheus.MustNewConstMetric(
				c.lossPercentage,
				prometheus.GaugeValue,
				float64(v.LossPercentage),
				v.Name,
				v.Address,
				c.instance,
			)
		}
		if v.RTTMilliseconds != -1 {
			ch <- prometheus.MustNewConstMetric(
				c.rtt,
				prometheus.GaugeValue,
				v.RTTMilliseconds,
				v.Name,
				v.Address,
				c.instance,
			)
		}
		if v.RTTDMilliseconds != -1 {
			ch <- prometheus.MustNewConstMetric(
				c.rttd,
				prometheus.GaugeValue,
				v.RTTDMilliseconds,
				v.Name,
				v.Address,
				c.instance,
			)
		}
	}
	return nil
}
