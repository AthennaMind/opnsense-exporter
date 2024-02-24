package collector

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "opnsense"

// CollectorInstance is the interface a service specific collectors must implement.
type CollectorInstance interface {
	Register(namespace, isntance string, log log.Logger)
	Name() string
	Describe(ch chan<- *prometheus.Desc)
	Update(client *opnsense.Client, ch chan<- prometheus.Metric) *opnsense.APICallError
}

// collectorInstances is a list of collectorInstances that will be registered
// from the init() function in each collector file
var collectorInstances []CollectorInstance

type Collector struct {
	instanceLabel string
	mutex         sync.RWMutex
	Client        *opnsense.Client
	log           log.Logger
	collectors    []CollectorInstance

	scrapes        prometheus.CounterVec
	endpointErrors prometheus.CounterVec
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
	return withoutCollectorInstance("arp_table")
}

// WithoutCronCollector Option
// removes the cron collector from the list of collectors
func WithoutCronCollector() Option {
	return withoutCollectorInstance("cron")
}

// New creates a new Collector instance.
func New(client *opnsense.Client, log log.Logger, instanceName string, options ...Option) (*Collector, error) {

	c := Collector{
		Client:        client,
		log:           log,
		instanceLabel: instanceName,
		collectors:    collectorInstances,
	}

	for _, option := range options {
		if err := option(&c); err != nil {
			return nil, errors.Join(err, fmt.Errorf("failed to apply option"))
		}
	}

	for _, collector := range c.collectors {
		collector.Register(namespace, instanceName, c.log)
	}

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

	prometheus.MustRegister(c.scrapes)
	prometheus.MustRegister(c.endpointErrors)

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

	for _, collector := range c.collectors {
		collector.Describe(ch)
	}
}

// Collect implements the prometheus.Collector interface.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var wg sync.WaitGroup
	wg.Add(len(c.collectors))

	for _, collector := range c.collectors {
		go func(coll CollectorInstance) {
			if err := coll.Update(c.Client, ch); err != nil {
				level.Error(c.log).Log(
					"msg", "failed to update",
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
