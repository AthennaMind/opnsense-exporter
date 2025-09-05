package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

// namespace is the prefix for all metrics.
const namespace = "opnsense"

// instanceLabelName is the label name for the current instance that is used
// to identify the instance in the metrics when there are
// multiple instances of the exporter running.
const instanceLabelName = "opnsense_instance"

const (
	ArpTableSubsystem   = "arp_table"
	GatewaysSubsystem   = "gateways"
	CronTableSubsystem  = "cron"
	WireguardSubsystem  = "wireguard"
	IPsecSubsystem      = "ipsec"
	UnboundDNSSubsystem = "unbound_dns"
	InterfacesSubsystem = "interfaces"
	ProtocolSubsystem   = "protocol"
	OpenVPNSubsystem    = "openvpn"
	ServicesSubsystem   = "services"
	FirewallSubsystem   = "firewall"
	FirmwareSubsystem   = "firmware"
)

// CollectorInstance is the interface a service specific collectors must implement.
type CollectorInstance interface {
	Register(namespace, isntance string, log *slog.Logger)
	Name() string
	Describe(ch chan<- *prometheus.Desc)
	Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError
}

// collectorInstances is a list of collectorInstances that will be registered
// from the init() function in each collector file
var collectorInstances []CollectorInstance

type Collector struct {
	Client *opnsense.Client
	mutex  sync.RWMutex
	log    *slog.Logger

	isUp                 prometheus.Gauge
	firewallHealthStatus prometheus.Gauge
	scrapes              prometheus.CounterVec
	endpointErrors       prometheus.CounterVec
	instanceLabel        string
	collectors           []CollectorInstance
}

type Option func(*Collector) error

// withoutCollectorInstance removes a collector by given name from the list of collectors
// that are registered from their init functions.
func withoutCollectorInstance(name string) Option {
	return func(o *Collector) error {
		for i, collector := range o.collectors {
			if collector.Name() == name {
				o.collectors = append(o.collectors[:i], o.collectors[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("collector %s not found", name)
	}
}

// WithoutArpTableCollector Option
// removes the arp_table collector from the list of collectors
func WithoutArpTableCollector() Option {
	return withoutCollectorInstance(ArpTableSubsystem)
}

// WithoutCronCollector Option
// removes the cron collector from the list of collectors
func WithoutCronCollector() Option {
	return withoutCollectorInstance(CronTableSubsystem)
}

// WithoutWireguardCollector Option
// removes the wireguard collector from the list of collectors
func WithoutWireguardCollector() Option {
	return withoutCollectorInstance(WireguardSubsystem)
}

// WithoutIPsecCollector Option
// removes the ipsec collector from the list of collectors
func WithoutIPsecCollector() Option {
	return withoutCollectorInstance(IPsecSubsystem)
}

// WithoutUnboundCollector Option
// removes the unbound_dns collector from the list of collectors
func WithoutUnboundCollector() Option {
	return withoutCollectorInstance(UnboundDNSSubsystem)
}

// WithoutFirewallCollector Option
// removes the firewall (pf) collector from the list of collectors
func WithoutFirewallCollector() Option {
	return withoutCollectorInstance(FirewallSubsystem)
}

// WithoutFirmwareCollector Option
// removes the firmware collector from the list of collectors
func WithoutFirmwareCollector() Option { return withoutCollectorInstance(FirmwareSubsystem) }

// WithoutOpenVPNCollector Option
// removes the openvpn collector from the list of collectors
func WithoutOpenVPNCollector() Option {
	return withoutCollectorInstance(OpenVPNSubsystem)
}

// New creates a new Collector instance.
func New(client *opnsense.Client, log *slog.Logger, instanceName string, options ...Option) (*Collector, error) {
	c := Collector{
		Client:        client,
		log:           log,
		instanceLabel: instanceName,
		collectors:    collectorInstances,
	}

	for _, option := range options {
		if err := option(&c); err != nil {
			return nil, errors.Join(err, fmt.Errorf("failed to apply collector option"))
		}
	}

	for _, collector := range c.collectors {
		collector.Register(namespace, instanceName, c.log)
	}

	c.isUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "up",
		Help:      "Was the last scrape of OPNsense successful. (1 = yes, 0 = no)",
		ConstLabels: prometheus.Labels{
			instanceLabelName: instanceName,
		},
	})

	c.firewallHealthStatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "firewall_status",
		Help:      "Status of the firewall reported by the system health check (1 = ok, 0 = errors)",
		ConstLabels: prometheus.Labels{
			instanceLabelName: instanceName,
		},
	})

	c.scrapes = *prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "exporter_scrapes_total",
		Help:      "Total number of times OPNsense was scraped for metrics.",
	}, []string{"opnsense_instance"})

	c.endpointErrors = *prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "exporter_endpoint_errors_total",
		Help:      "Total number of errors by endpoint returned by the OPNsense API during data fetching",
	}, []string{"endpoint", "opnsense_instance"})

	for _, metric := range []prometheus.Collector{c.isUp, c.scrapes, c.endpointErrors} {
		prometheus.MustRegister(metric)
	}

	c.scrapes.WithLabelValues(c.instanceLabel).Add(0)

	for _, path := range c.Client.Endpoints() {
		c.endpointErrors.WithLabelValues(string(path), c.instanceLabel).Add(0)
	}
	return &c, nil
}

// Describe implements the prometheus.Collector interface.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.scrapes.Describe(ch)
	c.endpointErrors.Describe(ch)
	c.isUp.Describe(ch)

	for _, collector := range c.collectors {
		collector.Describe(ch)
	}
}

func (c *Collector) collectHealthMetrics(ch chan<- prometheus.Metric) error {
	systemStatus, err := c.Client.HealthCheck()
	if err != nil {
		c.isUp.Set(0)
		c.isUp.Collect(ch)
		return err
	}

	if systemStatus.System.Status != opnsense.HealthCheckStatusOK &&
		systemStatus.Metadata.System.Status != opnsense.HealthCheckStatusOK_v25_1 {
		c.isUp.Set(0)
		c.isUp.Collect(ch)
		return nil
	}

	c.isUp.Set(1)
	c.firewallHealthStatus.Set(1)

	if systemStatus.Firewall.Status != opnsense.HealthCheckStatusOK &&
		systemStatus.Metadata.Firewall.Status != opnsense.HealthCheckStatusOK_v25_1 {
		c.firewallHealthStatus.Set(0)
	}

	c.isUp.Collect(ch)
	c.firewallHealthStatus.Collect(ch)
	return nil
}

// Collect implements the prometheus.Collector interface.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.collectHealthMetrics(ch); err != nil {
		c.log.Error(
			"failed to fetch system health status; skipping other metrics",
			"err", err,
		)
	}

	var wg sync.WaitGroup
	wg.Add(len(c.collectors))

	for _, collector := range c.collectors {
		go func(coll CollectorInstance) {
			if err := coll.Update(c.Client, ch); err != nil {
				c.log.Error(
					"failed to update",
					"component", "collector",
					"collector_name", coll.Name(),
					"err", err,
				)
				c.endpointErrors.WithLabelValues(err.Endpoint, c.instanceLabel).Inc()
			}
			wg.Done()
		}(collector)
	}
	wg.Wait()

	c.scrapes.WithLabelValues(c.instanceLabel).Inc()
	c.scrapes.Collect(ch)
	c.endpointErrors.Collect(ch)
}
