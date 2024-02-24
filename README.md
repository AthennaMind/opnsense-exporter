# OPNsense Prometheus Exporter

The missing OPNsense exporter for Prometheus

`Still under heavy development. The full metrics list is not yet implemented.`

## Table of Contents

- **[About](#about)**
- **[OPNsense User Permissions](#opnsense-user-permissions)**
- **[Development](#development)**
- **[Usage](#usage)**  
  - **[Docker](#docker)**
  - **[Docker Compose](#docker-compose)**
  - **[Systemd](#systemd)**
  - **[K8s](#k8s)**
- **[Configuration](#configuration)**  
  - **[SSL/TLS](#ssltls)**
  - **[All Options](#all-options)**
- **[Grafana Dashboard](#grafana-dashboard)**  

## About

Focusing specifically on OPNsense, this exporter provides metrics about OPNsense, the plugin ecosystem and the services running on the firewall. However, it's recommended to use it with `node_exporter`. You can combine the metrics from both exporters in Grafana and in your Alert System to create a dashboard that displays the full picture of your system.

While the `node_exporter` must be installed on the firewall itself, this exporter can be installed on any machine that has network access to the OPNsense API.

## Development

This guide is for osx and Linux.

### Create API key and secret in OPNsense

`SYSTEM>ACCESS>USERS>[user]>API KEYS`

[OPNsense Documentation](https://docs.opnsense.org/development/how-tos/api.html#creating-keys)

### Run the exporter locally

```bash
OPS_ADDRESS="ops.example.com" OPS_API_KEY=your-api-key OPS_API_SECRET=your-api-secret make local-run
curl http://localhost:8080/metrics
```

### Before PR

- Make sure to sync the vendor if the dependencies have changed.

```bash
make sync-vendor
```

- Make sure to run the tests and linters.

```bash
make test
make lint
```

## OPNsense user permissions

**TODO**

## Usage

**TODO**

### Docker 

To run the exporter using Docker, you can use the following command:

```bash
docker run -p 8080:8080 ghcr.io/AthennaMind/opnsense-exporter:latest \
      /opnsense-exporter \
      --log.level=debug \
      --log.format=json \
      --opnsense.protocol=https \
      --opnsense.address=ops.example.com \
      --opnsense.api-key=your-api-key \
      --opnsense.api-secret=your-api-secret \
      --exporter.instance-label=opnsense-eu1 \
      --web.listen-address=:8080 
```
### Docker Compose

**TODO**

### Systemd

**TODO**

### K8s

**TODO**

## Configuration

To configure where your OPNsense API is located, you can use the following flags:

- `--opnsense.protocol` - The protocol to use to connect to the OPNsense API. Can be either `http` or `https`.
- `--opnsense.address` - The hostname or IP address of the OPNsense API.
- `--opnsense.api-key` - The API key to use to connect to the OPNsense API.
- `--opnsense.api-secret` - The API secret to use to connect to the OPNsense API
- `--exporter.instance-label` - Label to use to identify the instance in every metric. If you have multiple instances of the exporter, you can differentiate them by using different value in this flag, that represents the instance of the target OPNsense.

### SSL/TLS

If you have your api served with self-signed certificates. You should add them to the system trust store.

If you want to disable TLS certificate verification, you can use the following flag:

- `--opnsense.insecure` - Disable TLS certificate verification. Defaults to `false`.

You can disable parts of the exporter using the following flags:

- `--exporter.disable-arp-table` - Disable the scraping of the ARP table. Defaults to `false`.
- `--exporter.disable-cron-table` - Disable the scraping of the cron table. Defaults to `false`.

You can disable the exporter metrics using the following flag:

- `--web.disable-exporter-metrics` - Exclude metrics about the exporter itself (promhttp_*, process_*, go_*). Defaults to `false`.

### All Options

```bash
Flags:
  -h, --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
      --log.level="info"         Log level. One of: [debug, info, warn, error]
      --log.format="logfmt"      Log format. One of: [logfmt, json]
      --web.telemetry-path="/metrics"  
                                 Path under which to expose metrics.
      --[no-]web.disable-exporter-metrics  
                                 Exclude metrics about the exporter itself (promhttp_*, process_*, go_*). ($OPNSENSE_EXPORTER_DISABLE_EXPORTER_METRICS)
      --runtime.gomaxprocs=2     The target number of CPUs that the Go runtime will run on (GOMAXPROCS) ($GOMAXPROCS)
      --exporter.instance-label=EXPORTER.INSTANCE-LABEL  
                                 Label to use to identify the instance in every metric. If you have multiple instances of the exporter, you can differentiate them by using different value in this flag, that represents the instance of the target OPNsense.
                                 ($OPNSENSE_EXPORTER_INSTANCE_LABEL)
      --[no-]exporter.disable-arp-table  
                                 Disable the scraping of the ARP table ($OPNSENSE_EXPORTER_DISABLE_ARP_TABLE)
      --[no-]exporter.disable-cron-table  
                                 Disable the scraping of the cron table ($OPNSENSE_EXPORTER_DISABLE_CRON_TABLE)
      --opnsense.protocol=OPNSENSE.PROTOCOL  
                                 Protocol to use to connect to OPNsense API. One of: [http, https] ($OPNSENSE_EXPORTER_OPS_PROTOCOL)
      --opnsense.address=OPNSENSE.ADDRESS  
                                 Hostname or IP address of OPNsense API ($OPNSENSE_EXPORTER_OPS_API)
      --opnsense.api-key=OPNSENSE.API-KEY  
                                 API key to use to connect to OPNsense API ($OPNSENSE_EXPORTER_OPS_API_KEY)
      --opnsense.api-secret=OPNSENSE.API-SECRET  
                                 API secret to use to connect to OPNsense API ($OPNSENSE_EXPORTER_OPS_API_SECRET)
      --[no-]opnsense.insecure   Disable TLS certificate verification ($OPNSENSE_EXPORTER_OPS_INSECURE)
      --[no-]web.systemd-socket  Use systemd socket activation listeners instead of port listeners (Linux only).
      --web.listen-address=:8080 ...  
                                 Addresses on which to expose metrics and web interface. Repeatable for multiple addresses.
      --web.config.file=""       [EXPERIMENTAL] Path to configuration file that can enable TLS or authentication. See: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md
```

## Grafana Dashboard

**TODO**
