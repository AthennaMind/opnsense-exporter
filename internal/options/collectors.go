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
	unboundCollectorDisabled = kingpin.Flag(
		"exporter.disable-unbound",
		"Disable the scraping of Unbound service",
	).Envar("OPNSENSE_EXPORTER_DISABLE_UNBOUND").Default("false").Bool()
	openVPNCollectorDisabled = kingpin.Flag(
		"exporter.disable-openvpn",
		"Disable the scraping of OpenVPN service",
	).Envar("OPNSENSE_EXPORTER_DISABLE_OPENVPN").Default("false").Bool()
	firewallCollectorDisabled = kingpin.Flag(
		"exporter.disable-firewall",
		"Disable the scraping of the firewall (pf) metrics",
	).Envar("OPNSENSE_EXPORTER_DISABLE_FIREWALL").Default("false").Bool()
)

// CollectorsDisableSwitch hold the enabled/disabled state of the collectors
type CollectorsDisableSwitch struct {
	ARP       bool
	Cron      bool
	Wireguard bool
	Unbound   bool
	OpenVPN   bool
	Firewall  bool
}

// CollectorsSwitches returns configured instances of CollectorsDisableSwitch
func CollectorsSwitches() CollectorsDisableSwitch {
	return CollectorsDisableSwitch{
		ARP:       !*arpTableCollectorDisabled,
		Cron:      !*cronTableCollectorDisabled,
		Wireguard: !*wireguardCollectorDisabled,
		Unbound:   !*unboundCollectorDisabled,
		OpenVPN:   !*openVPNCollectorDisabled,
		Firewall:  !*firewallCollectorDisabled,
	}
}
