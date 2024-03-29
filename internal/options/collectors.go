package options

import "github.com/alecthomas/kingpin/v2"

var (
	arpTableCollectorDisabled = kingpin.Flag(
		"exporter.disable-arp-table",
		"Disable the scraping of the ARP table",
	).Envar("OPNSENSE_EXPORTER_DISABLE_ARP_TABLE").Default("false").Bool()
	cronTableCollectorDisabled = kingpin.Flag(
		"exporter.disable-cron-table",
		"Disable the scraping of the cron table",
	).Envar("OPNSENSE_EXPORTER_DISABLE_CRON_TABLE").Default("false").Bool()
	wireguardCollectorDisabled = kingpin.Flag(
		"exporter.disable-wireguard",
		"Disable the scraping of Wireguard service",
	).Envar("OPNSENSE_EXPORTER_DISABLE_WIREGUARD").Default("false").Bool()
)

// Collectors holds the configuration for the collectors
type CollectorsSwitches struct {
	ARP       bool
	Cron      bool
	Wireguard bool
}

// Collectors returns the configuration for the collectors
func Collectors() CollectorsSwitches {
	return CollectorsSwitches{
		ARP:       !*arpTableCollectorDisabled,
		Cron:      !*cronTableCollectorDisabled,
		Wireguard: !*wireguardCollectorDisabled,
	}
}
