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
	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/exporter-toolkit/web"
)

var version = ""

func main() {
	options.Init()
	logger := promslog.New(options.PromLogConfig)

	runtime.GOMAXPROCS(*options.MaxProcs)

	logger.Info("starting opnsense-exporter", "version", version)
	logger.Info("settings Go MAXPROCS", "procs", runtime.GOMAXPROCS(0))

	opnsConfig, err := options.OPNSense()
	if err != nil {
		logger.Error("failed to assemble OPNsense configuration", "err", err)
		os.Exit(1)
	}

	opnsenseClient, err := opnsense.NewClient(
		*opnsConfig,
		version,
		logger,
	)
	if err != nil {
		logger.Error("opnsense client build failed", "err", err)
		os.Exit(1)
	}

	logger.Debug(fmt.Sprintf("OPNsense registered endpoints %s", opnsenseClient.Endpoints()))

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
		logger.Info("unbound collector disabled")
	}
	if !collectorsSwitches.Wireguard {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutWireguardCollector())
		logger.Info("wireguard collector disabled")
	}
	if !collectorsSwitches.IPsec {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutIPsecCollector())
		logger.Info("ipesc collector disabled")
	}
	if !collectorsSwitches.Cron {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutCronCollector())
		logger.Info("cron collector disabled")
	}
	if !collectorsSwitches.ARP {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutArpTableCollector())
		logger.Info("arp collector disabled")
	}
	if !collectorsSwitches.Firewall {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutFirewallCollector())
		logger.Info("firewall collector disabled")
	}
	if !collectorsSwitches.Firmware {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutFirmwareCollector())
		logger.Info("firmware collector disabled")
	}
	if !collectorsSwitches.OpenVPN {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutOpenVPNCollector())
		logger.Info("openvpn collector disabled")
	}

	collectorInstance, err := collector.New(&opnsenseClient, logger, *options.InstanceLabel, collectorOptionFuncs...)
	if err != nil {
		logger.Error("failed to construct the collecotr", "err", err)
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
			logger.Error("failed to construct landing page", "err", err)
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
			logger.Error("Error received from the HTTP server", "err", err)
			close(srvClose)
		}
	}()

	for {
		select {
		case <-term:
			logger.Info("Received SIGTERM, exiting gracefully...")
			os.Exit(0)
		case <-srvClose:
			os.Exit(1)
		}
	}
}
