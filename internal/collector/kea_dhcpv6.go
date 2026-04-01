package collector

import (
	"log/slog"
	"strconv"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

type keaDhcpv6Collector struct {
	log                      *slog.Logger
	lease_count              *prometheus.Desc
	lease_expiration         *prometheus.Desc
	lease_preferred_lifetime *prometheus.Desc
	leases_reserved          *prometheus.Desc
	lease_lifetime           *prometheus.Desc

	subsystem string
	instance  string
}

func init() {
	collectorInstances = append(collectorInstances, &keaDhcpv6Collector{
		subsystem: KeaDHCPv6Subsystem,
	})
}

func (c *keaDhcpv6Collector) Name() string { return c.subsystem }

func (c *keaDhcpv6Collector) Register(namespace string, instanceLabel string, log *slog.Logger) {
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
		[]string{"hostname", "ip_address", "prefix_len", "hwaddr", "duid", "if", "interface_name", "interface_description"},
	)
	c.lease_preferred_lifetime = buildPrometheusDesc(c.subsystem,
		"lease_preferred_lifetime",
		"Preferred lifetime of the lease",
		[]string{"hostname", "ip_address", "prefix_len", "hwaddr", "duid", "if", "interface_name", "interface_description"},
	)
	c.lease_lifetime = buildPrometheusDesc(c.subsystem,
		"lease_lifetime",
		"Lifetime for the lease",
		[]string{"hostname", "ip_address", "prefix_len", "hwaddr", "duid", "if", "interface_name", "interface_description"},
	)
	c.leases_reserved = buildPrometheusDesc(c.subsystem,
		"leases_reserved",
		"Number of reserved IP addresses",
		[]string{"interface_name"},
	)
}

func (c *keaDhcpv6Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.lease_count
	ch <- c.lease_expiration
	ch <- c.leases_reserved
	ch <- c.lease_lifetime
	ch <- c.lease_preferred_lifetime
}

func (c *keaDhcpv6Collector) Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError {
	keaDhcpv6Leases, err := client.FetchLeasesv6()
	if err != nil {
		return err
	}

	for _, lease := range keaDhcpv6Leases.Leases {
		ch <- prometheus.MustNewConstMetric(
			c.lease_lifetime,
			prometheus.GaugeValue,
			float64(lease.ValidLifetime),
			lease.Hostname,
			lease.Address,
			strconv.Itoa(lease.PrefixLength),
			lease.Hwaddr,
			lease.Duid,
			lease.InterfaceName,
			keaDhcpv6Leases.Interfaces[lease.InterfaceName].Name,
			keaDhcpv6Leases.Interfaces[lease.InterfaceName].Description,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.lease_expiration,
			prometheus.GaugeValue,
			float64(lease.Expiration),
			lease.Hostname,
			lease.Address,
			strconv.Itoa(lease.PrefixLength),
			lease.Hwaddr,
			lease.Duid,
			lease.InterfaceName,
			keaDhcpv6Leases.Interfaces[lease.InterfaceName].Name,
			keaDhcpv6Leases.Interfaces[lease.InterfaceName].Description,
			c.instance,
		)
		ch <- prometheus.MustNewConstMetric(
			c.lease_preferred_lifetime,
			prometheus.GaugeValue,
			float64(lease.PreferredLifetime),
			lease.Hostname,
			lease.Address,
			strconv.Itoa(lease.PrefixLength),
			lease.Hwaddr,
			lease.Duid,
			lease.InterfaceName,
			keaDhcpv6Leases.Interfaces[lease.InterfaceName].Name,
			keaDhcpv6Leases.Interfaces[lease.InterfaceName].Description,
			c.instance,
		)
	}

	for interfaceName, reservations := range keaDhcpv6Leases.ReservedLeaseCount {
		ch <- prometheus.MustNewConstMetric(
			c.leases_reserved,
			prometheus.GaugeValue,
			float64(reservations),
			interfaceName,
			c.instance,
		)
	}

	for interfaceName, activeLeases := range keaDhcpv6Leases.LeaseCount {
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
