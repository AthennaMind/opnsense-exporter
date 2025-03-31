package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

func buildPrometheusDesc(subsystem, name, help string, labels []string) *prometheus.Desc {
	if labels != nil {
		labels = append(labels, instanceLabelName)
	} else {
		labels = []string{instanceLabelName}
	}

	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, name),
		help,
		labels,
		nil,
	)
}

// parseStringToBool parses a string value to a bool value.
// The value is considered true if it is not "0".
func parseStringToBool(value string) bool {
	return value != "0"
}
