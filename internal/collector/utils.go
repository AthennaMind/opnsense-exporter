package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

const instanceLabelName = "opnsense_instance"

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
