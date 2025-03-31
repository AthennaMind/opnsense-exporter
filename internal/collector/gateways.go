package collector

import (
	"log/slog"
	"strconv"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type gatewaysCollector struct {
	log            *slog.Logger
	info           *prometheus.Desc
	monitor        *prometheus.Desc
	rtt            *prometheus.Desc
	rttd           *prometheus.Desc
	rttLow         *prometheus.Desc
	rttHigh        *prometheus.Desc
	lossPercentage *prometheus.Desc
	lossLow        *prometheus.Desc
	lossHigh       *prometheus.Desc
	interval       *prometheus.Desc
	period         *prometheus.Desc
	timeout        *prometheus.Desc
	status         *prometheus.Desc
	subsystem      string
	instance       string
}

func init() {
	collectorInstances = append(collectorInstances,
		&gatewaysCollector{
			subsystem: GatewaysSubsystem,
		},
	)
}

func (c *gatewaysCollector) Name() string {
	return c.subsystem
}

func (c *gatewaysCollector) Register(namespace, instanceLabel string, log *slog.Logger) {
	c.log = log
	c.instance = instanceLabel
	c.log.Debug("Registering collector", "collector", c.Name())

	c.info = buildPrometheusDesc(c.subsystem, "info",
		"Information of the gateway",
		[]string{"name", "description", "device", "protocol", "enabled", "weight", "interface", "upstream"},
	)
	c.monitor = buildPrometheusDesc(
		c.subsystem, "monitor_info",
		"Gateway monitoring configuration",
		[]string{"name", "enabled", "no_route", "address"},
	)
	c.rtt = buildPrometheusDesc(
		c.subsystem, "rtt_milliseconds",
		"RTT is the average (mean) of the round trip time in milliseconds by name and address",
		[]string{"name", "address"},
	)
	c.rttd = buildPrometheusDesc(
		c.subsystem, "rttd_milliseconds",
		"RTTd is the standard deviation of the round trip time in milliseconds by name and address",
		[]string{"name", "address"},
	)
	c.rttLow = buildPrometheusDesc(
		c.subsystem, "rtt_low_milliseconds",
		"Gateway low latency threshold",
		[]string{"name", "address"},
	)
	c.rttHigh = buildPrometheusDesc(
		c.subsystem, "rtt_high_milliseconds",
		"Gateway high latency threshold",
		[]string{"name", "address"},
	)
	c.lossPercentage = buildPrometheusDesc(
		c.subsystem, "loss_percentage",
		"The current gateway loss percentage by name and address",
		[]string{"name", "address"},
	)
	c.lossLow = buildPrometheusDesc(
		c.subsystem, "loss_low_percentage",
		"Gateway low packet loss threshold",
		[]string{"name", "address"},
	)
	c.lossHigh = buildPrometheusDesc(
		c.subsystem, "loss_high_percentage",
		"Gateway high packet loss threshold",
		[]string{"name", "address"},
	)
	c.interval = buildPrometheusDesc(
		c.subsystem, "probe_interval_seconds",
		"Gateway probe interval",
		[]string{"name", "address"},
	)
	c.period = buildPrometheusDesc(
		c.subsystem, "probe_period_seconds",
		"Gateway probe period",
		[]string{"name", "address"},
	)
	c.timeout = buildPrometheusDesc(
		c.subsystem, "probe_timeout_seconds",
		"Gateway probe timeout",
		[]string{"name", "address"},
	)
	c.status = buildPrometheusDesc(c.subsystem, "status",
		"Status of the gateway by name and address (0 = Offline, 1 = Online, 2 = Unknown, 3 = Pending)",
		[]string{"name", "address", "default_gateway"},
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
		monitorEnabledFloat := 1.0
		if !v.MonitorEnabled {
			monitorEnabledFloat = 0.0
		}
		interfaceEnabledFloat := 1.0
		if !v.Enabled {
			interfaceEnabledFloat = 0.0
		}

		ch <- prometheus.MustNewConstMetric(
			c.info,
			prometheus.GaugeValue,
			interfaceEnabledFloat,
			v.Name,
			v.Description,
			v.Interface,
			v.IPProtocol,
			strconv.FormatBool(v.Enabled),
			v.Weight,
			v.HardwareInterface,
			strconv.FormatBool(v.Upstream),
			c.instance,
		)
		if v.Enabled {
			ch <- prometheus.MustNewConstMetric(
				c.monitor,
				prometheus.GaugeValue,
				monitorEnabledFloat,
				v.Name,
				strconv.FormatBool(v.MonitorEnabled),
				strconv.FormatBool(v.MonitorNoRoute),
				v.Monitor,
				c.instance,
			)
			if v.MonitorEnabled {
				ch <- prometheus.MustNewConstMetric(
					c.rtt,
					prometheus.GaugeValue,
					v.Delay,
					v.Name,
					v.Monitor,
					c.instance,
				)
				ch <- prometheus.MustNewConstMetric(
					c.rttd,
					prometheus.GaugeValue,
					v.StdDev,
					v.Name,
					v.Monitor,
					c.instance,
				)
				f64, _ := strconv.ParseFloat(v.LatencyLow, 64)
				ch <- prometheus.MustNewConstMetric(
					c.rttLow,
					prometheus.GaugeValue,
					f64,
					v.Name,
					v.Monitor,
					c.instance,
				)
				f64, _ = strconv.ParseFloat(v.LatencyHigh, 64)
				ch <- prometheus.MustNewConstMetric(
					c.rttHigh,
					prometheus.GaugeValue,
					f64,
					v.Name,
					v.Monitor,
					c.instance,
				)
				ch <- prometheus.MustNewConstMetric(
					c.lossPercentage,
					prometheus.GaugeValue,
					v.Loss,
					v.Name,
					v.Monitor,
					c.instance,
				)
				f64, _ = strconv.ParseFloat(v.LossLow, 64)
				ch <- prometheus.MustNewConstMetric(
					c.lossLow,
					prometheus.GaugeValue,
					f64,
					v.Name,
					v.Monitor,
					c.instance,
				)
				f64, _ = strconv.ParseFloat(v.LossHigh, 64)
				ch <- prometheus.MustNewConstMetric(
					c.lossHigh,
					prometheus.GaugeValue,
					f64,
					v.Name,
					v.Monitor,
					c.instance,
				)
				f64, _ = strconv.ParseFloat(v.Interval, 64)
				ch <- prometheus.MustNewConstMetric(
					c.interval,
					prometheus.GaugeValue,
					f64,
					v.Name,
					v.Monitor,
					c.instance,
				)
				f64, _ = strconv.ParseFloat(v.LossInterval, 64)
				ch <- prometheus.MustNewConstMetric(
					c.timeout,
					prometheus.GaugeValue,
					f64,
					v.Name,
					v.Monitor,
					c.instance,
				)
				ch <- prometheus.MustNewConstMetric(
					c.status,
					prometheus.GaugeValue,
					float64(v.Status),
					v.Name,
					v.Monitor,
					strconv.FormatBool(v.DefaultGateway),
					c.instance,
				)
			}
		}
	}
	return nil
}
