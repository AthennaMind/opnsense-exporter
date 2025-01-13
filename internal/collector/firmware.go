package collector

import (
	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type firmwareCollector struct {
	log log.Logger

	lastCheck          *prometheus.Desc
	needsReboot        *prometheus.Desc
	newPackages        *prometheus.Desc
	osVersion          *prometheus.Desc
	productAbi         *prometheus.Desc
	productId          *prometheus.Desc
	productVersion     *prometheus.Desc
	upgradePackages    *prometheus.Desc
	upgradeNeedsReboot *prometheus.Desc

	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &firmwareCollector{
		subsystem: FirmwareSubsystem,
	})
}

func (c *firmwareCollector) Name() string {
	return c.subsystem
}

func (c *firmwareCollector) Register(namespace, instanceLabel string, log log.Logger) {
	c.log = log
	c.instance = instanceLabel

	level.Debug(c.log).
		Log("msg", "Registering collector", "collector", c.Name())

	c.lastCheck = buildPrometheusDesc(c.subsystem, "last_check",
		"last check for upgrade", []string{"last_check"})

	c.needsReboot = buildPrometheusDesc(c.subsystem, "needs_reboot",
		"opnsense would like to be rebooted", []string{"needs_reboot"})

	c.newPackages = buildPrometheusDesc(c.subsystem, "new_packages",
		"new packages", []string{"new_packages"})

	c.osVersion = buildPrometheusDesc(c.subsystem, "os_version",
		"Version of this opnSense", []string{"os_version"})

	c.productAbi = buildPrometheusDesc(c.subsystem, "product_abi",
		"Product ABI of this opnSense", []string{"product_abi"})

	c.productId = buildPrometheusDesc(c.subsystem, "product_id",
		"Product ID of this opnSense", []string{"product_id"})

	c.productVersion = buildPrometheusDesc(c.subsystem, "product_version",
		"Product Version of this opnSense", []string{"product_version"})

	c.upgradePackages = buildPrometheusDesc(c.subsystem, "upgrade_packages",
		"upgrade packages", []string{"upgrade_packages"})

	c.upgradeNeedsReboot = buildPrometheusDesc(c.subsystem, "upgrade_needs_reboot",
		"upgrade involves reboot", []string{"upgrade_needs_reboot"})
}

func (c *firmwareCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.needsReboot
	ch <- c.newPackages
	ch <- c.lastCheck
	ch <- c.osVersion
	ch <- c.productAbi
	ch <- c.productId
	ch <- c.productVersion
	ch <- c.upgradePackages
	ch <- c.upgradeNeedsReboot
}

func (c *firmwareCollector) update(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(
		desc, valueType, value, labelValues...,
	)
}

func (c *firmwareCollector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	data, err := client.FetchFirmwareStatus()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.needsReboot, prometheus.GaugeValue, float64(data.NeedsReboot), strconv.Itoa(data.NeedsReboot), c.instance)
	ch <- prometheus.MustNewConstMetric(c.newPackages, prometheus.GaugeValue, float64(data.NewPackages), strconv.Itoa(data.NewPackages), c.instance)
	ch <- prometheus.MustNewConstMetric(c.lastCheck, prometheus.GaugeValue, float64(1), data.LastCheck, c.instance)
	ch <- prometheus.MustNewConstMetric(c.osVersion, prometheus.GaugeValue, float64(1), data.OsVersion, c.instance)
	ch <- prometheus.MustNewConstMetric(c.productAbi, prometheus.GaugeValue, float64(1), data.ProductABI, c.instance)
	ch <- prometheus.MustNewConstMetric(c.productId, prometheus.GaugeValue, float64(1), data.ProductId, c.instance)
	ch <- prometheus.MustNewConstMetric(c.productVersion, prometheus.GaugeValue, float64(1), data.ProductVersion, c.instance)
	ch <- prometheus.MustNewConstMetric(c.upgradePackages, prometheus.GaugeValue, float64(data.UpgradePackages), strconv.Itoa(data.UpgradePackages), c.instance)
	ch <- prometheus.MustNewConstMetric(c.upgradeNeedsReboot, prometheus.GaugeValue, float64(data.UpgradeNeedsReboot), strconv.Itoa(data.UpgradeNeedsReboot), c.instance)

	return nil
}
