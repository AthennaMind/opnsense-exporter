package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"github.com/st3ga/opnsense-exporter/internal/collector"
	"github.com/st3ga/opnsense-exporter/opnsense"
)

var version = ""

func main() {
	var (
		logLevel = kingpin.Flag(
			"log.level",
			"Log level. One of: [debug, info, warn, error]").
			Default("info").
			String()
		logFormat = kingpin.Flag(
			"log.format",
			"Log format. One of: [logfmt, json]").
			Default("logfmt").
			String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		disableExporterMetrics = kingpin.Flag(
			"web.disable-exporter-metrics",
			"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
		).Envar("OPNSENSE_EXPORTER_DISABLE_EXPORTER_METRICS").Bool()
		maxProcs = kingpin.Flag(
			"runtime.gomaxprocs",
			"The target number of CPUs that the Go runtime will run on (GOMAXPROCS)",
		).Envar("GOMAXPROCS").Default("2").Int()
		instanceLabel = kingpin.Flag(
			"exporter.instance-label",
			"Label to use to identify the instance in every metric. "+
				"If you have multiple instances of the exporter, you can differentiate them by using "+
				"different value in this flag, that represents the instance of the target OPNsense.",
		).Envar("OPNSENSE_EXPORTER_INSTANCE_LABEL").Required().String()
		arpTableCollectorDisabled = kingpin.Flag(
			"exporter.disable-arp-table",
			"Disable the scraping of the ARP table",
		).Envar("OPNSENSE_EXPORTER_DISABLE_ARP_TABLE").Default("false").Bool()
		cronTableCollectorDisabled = kingpin.Flag(
			"exporter.disable-cron-table",
			"Disable the scraping of the cron table",
		).Envar("OPNSENSE_EXPORTER_DISABLE_CRON_TABLE").Default("false").Bool()
		opnsenseProtocol = kingpin.Flag(
			"opnsense.protocol",
			"Protocol to use to connect to OPNsense API. One of: [http, https]",
		).Envar("OPNSENSE_EXPORTER_OPS_PROTOCOL").Required().String()
		opnsenseAPI = kingpin.Flag(
			"opnsense.address",
			"Hostname or IP address of OPNsense API",
		).Envar("OPNSENSE_EXPORTER_OPS_API").Required().String()
		opnsenseAPIKey = kingpin.Flag(
			"opnsense.api-key",
			"API key to use to connect to OPNsense API",
		).Envar("OPNSENSE_EXPORTER_OPS_API_KEY").Required().String()
		opnsenseAPISecret = kingpin.Flag(
			"opnsense.api-secret",
			"API secret to use to connect to OPNsense API",
		).Envar("OPNSENSE_EXPORTER_OPS_API_SECRET").Required().String()
		opnsenseInsecure = kingpin.Flag(
			"opnsense.insecure",
			"Disable TLS certificate verification",
		).Envar("OPNSENSE_EXPORTER_OPS_INSECURE").Default("false").Bool()

		webConfig = kingpinflag.AddFlags(kingpin.CommandLine, ":8080")
	)

	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	promlogConfig := &promlog.Config{
		Level:  &promlog.AllowedLevel{},
		Format: &promlog.AllowedFormat{},
	}
	promlogConfig.Level.Set(*logLevel)
	promlogConfig.Format.Set(*logFormat)

	logger := promlog.New(promlogConfig)

	level.Info(logger).
		Log("msg", "Starting opnsense-exporter", "version", version)

	runtime.GOMAXPROCS(*maxProcs)

	level.Debug(logger).
		Log("msg", "settings Go MAXPROCS", "procs", runtime.GOMAXPROCS(0))

	opnsenseClient, err := opnsense.NewClient(
		*opnsenseProtocol,
		*opnsenseAPI,
		*opnsenseAPIKey,
		*opnsenseAPISecret,
		version,
		*opnsenseInsecure,
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

	r := prometheus.NewRegistry()

	if !*disableExporterMetrics {
		r.MustRegister(
			promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
		)
		r.MustRegister(promcollectors.NewGoCollector())
	}

	collectorOptionFuncs := []collector.Option{}

	if *arpTableCollectorDisabled {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutArpTableCollector())
	}

	if *cronTableCollectorDisabled {
		collectorOptionFuncs = append(collectorOptionFuncs, collector.WithoutCronCollector())
	}

	collectorInstance, err := collector.New(&opnsenseClient, logger, *instanceLabel, collectorOptionFuncs...)

	if err != nil {
		level.Error(logger).
			Log("msg", "failed to construct the collecotr", "err", err)
		os.Exit(1)
	}

	r.MustRegister(collectorInstance)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	http.Handle(*metricsPath, handler)

	if *metricsPath != "/" && *metricsPath != "" {
		landingConfig := web.LandingConfig{
			Name:        "OPNsense Exporter",
			Description: "Prometheus OPNsense Firewall Exporter",
			Version:     version,
			Links: []web.LandingLinks{
				{
					Address: *metricsPath,
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
		if err := web.ListenAndServe(srv, webConfig, logger); err != nil {
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
