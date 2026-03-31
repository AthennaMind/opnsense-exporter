package collector

import (
	"log/slog"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type keaDhcpv4Collector struct {
	log              *slog.Logger
	lease_count      *prometheus.Desc
	lease_expiration *prometheus.Desc
	leases_reserved  *prometheus.Desc
	lease_lifetime   *prometheus.Desc

	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &keaDhcpv4Collector{
		subsystem: KeaDHCPv4Subsystem,
	})
}

func (c *keaDhcpv4Collector) Name() string { return c.subsystem }

func (c *keaDhcpv4Collector) Register(namespace string, instanceLabel string, log *slog.Logger) {
	c.log = log
	c.instance = instanceLabel

	c.log.Debug("Registering collector", "collector", c.Name())

	c.lease_count = buildPrometheusDesc(c.subsystem,
		"lease_count",
		"Number of leases for the interface",
		[]string{"interface_name"},
	)
	c.lease_expiration = buildPrometheusDesc(c.subsystem,
		"lease_expiration",
		"Time the lease expires in seconds",
		[]string{"hostname", "ip_address", "mac", "mac_info", "client_id", "if", "interface_name", "interface_description"},
	)
	c.lease_lifetime = buildPrometheusDesc(c.subsystem,
		"lease_lifetime",
		"Lifetime for the lease",
		[]string{"hostname", "ip_address", "mac", "mac_info", "client_id", "if", "interface_name", "interface_description"},
	)
	c.leases_reserved = buildPrometheusDesc(c.subsystem,
		"leases_reserved",
		"Number of reserved IP addresses",
		[]string{"interface_name"},
	)
}

func (c *keaDhcpv4Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.lease_count
	ch <- c.lease_expiration
	ch <- c.leases_reserved
	ch <- c.lease_lifetime
}

func (c *keaDhcpv4Collector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	keaDhcpv4Leases, err := client.FetchLeasesv4()
	if err != nil {
		return err
	}

	for _, lease := range keaDhcpv4Leases.Leases {
		ch <- prometheus.MustNewConstMetric(
			c.lease_lifetime,
			prometheus.GaugeValue,
			float64(lease.ValidLifetime),
			lease.Hostname,
			lease.Address,
			lease.Mac,
			lease.MacInfo,
			lease.ClientId,
			lease.InterfaceName,
			keaDhcpv4Leases.Interfaces[lease.InterfaceName].Name,
			keaDhcpv4Leases.Interfaces[lease.InterfaceName].Description,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.lease_expiration,
			prometheus.GaugeValue,
			float64(lease.Expiration),
			lease.Hostname,
			lease.Address,
			lease.Mac,
			lease.MacInfo,
			lease.ClientId,
			lease.InterfaceName,
			keaDhcpv4Leases.Interfaces[lease.InterfaceName].Name,
			keaDhcpv4Leases.Interfaces[lease.InterfaceName].Description,
			c.instance,
		)
	}

	for interfaceName, reservations := range keaDhcpv4Leases.ReservedLeaseCount {
		ch <- prometheus.MustNewConstMetric(
			c.leases_reserved,
			prometheus.GaugeValue,
			float64(reservations),
			interfaceName,
			c.instance,
		)
	}

	for interfaceName, activeLeases := range keaDhcpv4Leases.LeaseCount {
		ch <- prometheus.MustNewConstMetric(
			c.lease_count,
			prometheus.GaugeValue,
			float64(activeLeases),
			interfaceName,
			c.instance,
		)
	}

	return nil
}
