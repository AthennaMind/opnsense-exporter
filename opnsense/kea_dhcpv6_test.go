package opnsense

import (
	"strings"
	"testing"
)

func TestParseKeaDHCPv6Leases(t *testing.T) {
	tests := []struct {
		name     string
		row      KeaDhcpv6LeasesResponse
		expected KeaDhcpv6Leases
	}{{
		name:     "no leases",
		row:      KeaDhcpv6LeasesResponse{},
		expected: KeaDhcpv6Leases{},
	}, {
		name: "1 lease, 1 interface",
		row: KeaDhcpv6LeasesResponse{
			Total:    1,
			RowCount: 1,
			Current:  1,
			Rows: []KeaDhcpv6LeasesRow{{
				If:                   "tst1",
				Address:              "1.2.3.4",
				Hwaddr:               "01:23:45:67:89:ab",
				Duid:                 "01:23:45:67:89:ab",
				ValidLifetime:        "86400",
				Expiration:           "86400",
				SubnetId:             "1",
				FqdnForward:          "0",
				FqdnReceived:         "0",
				Hostname:             "test",
				State:                "0",
				UserContext:          "",
				PoolId:               "0",
				InterfaceDescription: "Test Interface",
				InterfaceName:        "opt1",
				PreferredLifetime:    "5400",
				IsReserved:           "",
				PrefixLength:         "128",
			}},
		},
		expected: KeaDhcpv6Leases{
			Leases: []KeaDhcpv6Lease{{
				Address:           "1.2.3.4",
				Hwaddr:            "01:23:45:67:89:ab",
				Duid:              "01:23:45:67:89:ab",
				ValidLifetime:     86400,
				Expiration:        86400,
				PreferredLifetime: 5400,
				Hostname:          "test",
				InterfaceName:     "opt1",
				PrefixLength:      128,
			}},
			LeaseCount: map[string]int{
				"opt1": 1,
			},
			Interfaces: map[string]KeaDhcpV6InterfaceInfo{
				"opt1": {
					Name:        "tst1",
					Description: "Test Interface",
				},
			},
		},
	}, {
		name: "1 lease, 1 interface, 1 reservation",
		row: KeaDhcpv6LeasesResponse{
			Total:    1,
			RowCount: 1,
			Current:  1,
			Rows: []KeaDhcpv6LeasesRow{{
				If:                   "tst1",
				Address:              "2008:db8:acad:dead:beef:cafe:7001:10af",
				Hwaddr:               "01:23:45:67:89:ab",
				Duid:                 "00:01:00:01:31:5f:2f:1a:01:23:45:67:89:ab",
				ValidLifetime:        "86400",
				Expiration:           "86400",
				SubnetId:             "1",
				FqdnForward:          "0",
				FqdnReceived:         "0",
				Hostname:             "test",
				State:                "0",
				UserContext:          "",
				PoolId:               "0",
				InterfaceDescription: "Test Interface",
				InterfaceName:        "opt1",
				IsReserved:           "mac",
				PreferredLifetime:    "5400",
				PrefixLength:         "128",
			}},
		},
		expected: KeaDhcpv6Leases{
			Leases: []KeaDhcpv6Lease{{
				Address:           "2008:db8:acad:dead:beef:cafe:7001:10af",
				Hwaddr:            "01:23:45:67:89:ab",
				Duid:              "00:01:00:01:31:5f:2f:1a:01:23:45:67:89:ab",
				ValidLifetime:     86400,
				Expiration:        86400,
				PreferredLifetime: 5400,
				Hostname:          "test",
				InterfaceName:     "opt1",
				PrefixLength:      128,
			}},
			LeaseCount: map[string]int{
				"opt1": 1,
			},
			ReservedLeaseCount: map[string]int{
				"opt1": 1,
			},
			Interfaces: map[string]KeaDhcpV6InterfaceInfo{
				"opt1": {
					Name:        "tst1",
					Description: "Test Interface",
				},
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := parseDHCPv6Leases(tt.row)
			if err != nil {
				t.Error(err)
			}

			// Make sure correct amount of interfaces come back
			if len(data.Interfaces) != len(tt.expected.Interfaces) {
				t.Errorf("unexpected number of interfaces in response, expected %d, got %d\n", len(tt.expected.Interfaces), len(data.Interfaces))
			}
			if len(data.Leases) != len(tt.expected.Leases) {
				t.Errorf("unexpected number of leases in response, expected %d, got %d\n", len(tt.expected.Interfaces), len(data.Interfaces))
			}
			if len(data.LeaseCount) != len(tt.expected.LeaseCount) {
				t.Errorf("unexpected number of interfaces with leases in response, expected %d, got %d\n", len(tt.expected.Interfaces), len(data.Interfaces))
			}
			if len(data.ReservedLeaseCount) != len(tt.expected.ReservedLeaseCount) {
				t.Errorf("unexpected number of interfaces with reservations in response, expected %d, got %d\n", len(tt.expected.Interfaces), len(data.Interfaces))
			}

			// Verify the leases come back as expected
			for idx, lease := range data.Leases {
				if strings.Compare(lease.InterfaceName, tt.expected.Leases[idx].InterfaceName) != 0 {
					t.Errorf("unexpected interface name: %s, expected %s\n", lease.InterfaceName, tt.expected.Leases[idx].InterfaceName)
				}
				if strings.Compare(lease.Hwaddr, tt.expected.Leases[idx].Hwaddr) != 0 {
					t.Errorf("unexpected MAC: %s, expected %s\n", lease.Hwaddr, tt.expected.Leases[idx].Hwaddr)
				}
				if strings.Compare(lease.Hostname, tt.expected.Leases[idx].Hostname) != 0 {
					t.Errorf("unexpected hostname: %s, expected %s\n", lease.Hostname, tt.expected.Leases[idx].Hostname)
				}
				if strings.Compare(lease.Duid, tt.expected.Leases[idx].Duid) != 0 {
					t.Errorf("unexpected duid: %s, expected %s\n", lease.Duid, tt.expected.Leases[idx].Duid)
				}
				if strings.Compare(lease.Hostname, tt.expected.Leases[idx].Hostname) != 0 {
					t.Errorf("unexpected hostname: %s, expected %s\n", lease.Hostname, tt.expected.Leases[idx].Hostname)
				}
				if strings.Compare(lease.Address, tt.expected.Leases[idx].Address) != 0 {
					t.Errorf("unexpected address: %s, expected %s\n", lease.Address, tt.expected.Leases[idx].Address)
				}
				if lease.Expiration != tt.expected.Leases[idx].Expiration {
					t.Errorf("unexpected expiration: %d, expected %d\n", lease.Expiration, tt.expected.Leases[idx].Expiration)
				}
				if lease.PreferredLifetime != tt.expected.Leases[idx].PreferredLifetime {
					t.Errorf("unexpected preferred lifetime: %d, expected %d\n", lease.PreferredLifetime, tt.expected.Leases[idx].PreferredLifetime)
				}
				if lease.PrefixLength != tt.expected.Leases[idx].PrefixLength {
					t.Errorf("unexpected prefix length: %d, expected %d\n", lease.PrefixLength, tt.expected.Leases[idx].PrefixLength)
				}
			}

			// Verify the reservations come back correct
			for ifName, reservation := range data.ReservedLeaseCount {
				if reservation != tt.expected.ReservedLeaseCount[ifName] {
					t.Errorf("unexpected reservations for %s, expected %d, got %d\n", ifName, reservation, tt.expected.ReservedLeaseCount[ifName])
				}
			}

			// Verify the leases come back correct
			for ifName, leases := range data.LeaseCount {
				if leases != tt.expected.LeaseCount[ifName] {
					t.Errorf("unexpected current leases for %s, expected %d, got %d\n", ifName, leases, tt.expected.LeaseCount[ifName])
				}
			}
		})
	}
}
