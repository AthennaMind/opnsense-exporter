package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/AthennaMind/opnsense-exporter/internal/collector"
	"github.com/AthennaMind/opnsense-exporter/internal/options"
	"github.com/AthennaMind/opnsense-exporter/opnsense"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/web"
)

var version = ""

func main() {
	options.Init()

	logger, err := options.Logger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating logger: %v\n", err)
		os.Exit(1)
	}

	level.Info(logger).
		Log("msg", "starting opnsense-exporter", "version", version)

	runtime.GOMAXPROCS(*options.MaxProcs)

	level.Debug(logger).
		Log("msg", "settings Go MAXPROCS", "procs", runtime.GOMAXPROCS(0))

	opnsConfig, err := options.OPNSense()
	if err != nil {
		level.Error(logger).
			Log("msg", "failed to assemble OPNsense configuration", "err", err)
		os.Exit(1)
	}

	opnsenseClient, err := opnsense.NewClient(
		*opnsConfig,
		version,
		logger,
	)
	if err != nil {
		level.Error(logger).
			Log("msg", "opnsense client build failed", "err", err)
		os.Exit(1)
	}

	level.Debug(logger).Log(
		"msg", fmt.Sprintf("OPNsense registered endpoints %s", opnsenseClient.Endpoints()),
	)

	registry := prometheus.NewRegistry()

	if !*options.DisableExporterMetrics {
		registry.MustRegister(
			promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
		)
		registry.MustRegister(promcollectors.NewGoCollector())
	}

	collectorsSwitches := options.CollectorsSwitches()
	collectorOptionFuncs := []collector.Option{}

	if !collectorsSwitches.Unbound {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutUnboundCollector())
		level.Info(logger).Log("msg", "unbound collector disabled")
	}
	if !collectorsSwitches.Wireguard {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutWireguardCollector())
		level.Info(logger).Log("msg", "wireguard collector disabled")
	}
	if !collectorsSwitches.Cron {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutCronCollector())
		level.Info(logger).Log("msg", "cron collector disabled")
	}
	if !collectorsSwitches.ARP {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutArpTableCollector())
		level.Info(logger).Log("msg", "arp collector disabled")
	}
	if !collectorsSwitches.Firewall {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutFirewallCollector())
		level.Info(logger).Log("msg", "firewall collector disabled")
	}

	collectorInstance, err := collector.New(&opnsenseClient, logger, *options.InstanceLabel, collectorOptionFuncs...)
	if err != nil {
		level.Error(logger).
			Log("msg", "failed to construct the collecotr", "err", err)
		os.Exit(1)
	}

	registry.MustRegister(collectorInstance)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.Handle(*options.MetricsPath, handler)

	if *options.MetricsPath != "/" && *options.MetricsPath != "" {
		landingConfig := web.LandingConfig{
			Name:        "OPNsense Exporter",
			Description: "Prometheus OPNsense Firewall Exporter",
			Version:     version,
			Links: []web.LandingLinks{
				{
					Address: *options.MetricsPath,
					Text:    "Metrics",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			level.Error(logger).Log("err", err)
			os.Exit(1)
		}
		http.Handle("/", landingPage)
	}

	term := make(chan os.Signal, 1)
	srvClose := make(chan struct{})
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{}
	go func() {
		if err := web.ListenAndServe(srv, options.WebConfig, logger); err != nil {
			level.Error(logger).
				Log("msg", "Error received from the HTTP server", "err", err)
			close(srvClose)
		}
	}()

	for {
		select {
		case <-term:
			level.Info(logger).
				Log("msg", "Received SIGTERM, exiting gracefully...")
			os.Exit(0)
		case <-srvClose:
			os.Exit(1)
		}
	}
}
