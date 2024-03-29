package options

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

var (
	MetricsPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()
	DisableExporterMetrics = kingpin.Flag(
		"web.disable-exporter-metrics",
		"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
	).Envar("OPNSENSE_EXPORTER_DISABLE_EXPORTER_METRICS").Bool()
	MaxProcs = kingpin.Flag(
		"runtime.gomaxprocs",
		"The target number of CPUs that the Go runtime will run on (GOMAXPROCS)",
	).Envar("GOMAXPROCS").Default("2").Int()
	InstanceLabel = kingpin.Flag(
		"exporter.instance-label",
		"Label to use to identify the instance in every metric. "+
			"If you have multiple instances of the exporter, you can differentiate them by using "+
			"different value in this flag, that represents the instance of the target OPNsense.",
	).Envar("OPNSENSE_EXPORTER_INSTANCE_LABEL").Required().String()

	WebConfig = kingpinflag.AddFlags(kingpin.CommandLine, ":8080")
)
