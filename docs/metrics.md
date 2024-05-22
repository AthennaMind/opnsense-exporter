## OPNsense Exporter Metrics List

This table represent each metric and it's labels, the subsystem that it belongs, its description and how to disable it. The opnsense_instance label is applied to all metrics.


### General 

| Metric Name | Type | Labels | Subsystem | Description | Disable Flag |
| --- | --- | --- | --- | --- | --- |
opnsense_up | Gauge | opnsense_instance | n/a | The current status of OPNsense (1 = up, 0 = down) | n/a |
opnsense_firewall_status | Gauge | opnsense_instance | n/a | Status of the firewall reported by `api/core/system/status` ( 1 = ok, 0 = errors) | n/a |
opnsense_exporter_scrapes_total | Counter | n/a | n/a | Total number of scrapes by the OPNsense exporter | n/a |
opnsense_exporter_endpoint_errors_total | Counter | endpoint | n/a | Total number of errors by endpoint returned by the OPNsense API during data fetching | n/a |
opnsense_cron_job_status | Gauge | command, description, origin, schedule | Cron Table | Cron job status by name and description (1 = enabled, 0 = disabled) | --exporter.disable-cron-table |


### Services 

| Metric Name | Type | Labels | Subsystem | Description | Disable Flag |
| --- | --- | --- | --- | --- | --- |
opnsense_services_running_total | Gauge | n/a | Services | Total number of running services | n/a |
opnsense_services_stopped_total | Gauge | n/a | Services | Total number of stopped services | n/a |
opnsense_services_status | Gauge | name, description | Services | Service status by name and description (1 = running, 0 = stopped) | n/a |


### Interfaces

| Metric Name | Type | Labels | Subsystem | Description | Disable Flag |
| --- | --- | --- | --- | --- | --- |
opnsense_interfaces_mtu_bytes | Gauge | interface, device, type | Interfaces | MTU of the interface by interface name and device | n/a |
opnsense_interfaces_received_bytes_total | Counter | interface, device, type | Interfaces | Total number of received bytes on this interface by interface name and device | n/a |
opnsense_interfaces_transmitted_bytes_total | Counter | interface, device, type | Interfaces | Total number of transmitted bytes on this interface by interface name and device | n/a |
opnsense_interfaces_received_multicasts_total | Counter | interface, device, type | Interfaces | Total number of received multicast packets on this interface by interface name and device | n/a |
opnsense_interfaces_transmitted_multicasts_total | Counter | interface, device, type | Interfaces | Total number of transmitted multicast packets on this interface by interface name and device | n/a |
opnsense_interfaces_input_errors_total | Counter | interface, device, type | Interfaces | Input errors on this interface by interface name and device | n/a |
opnsense_interfaces_output_errors_total | Counter | interface, device, type | Interfaces | Output errors on this interface by interface name and device | n/a |
opnsense_interfaces_collisions_total | Counter | interface, device, type | Interfaces | Collisions on this interface by interface name and device | n/a |

![interfaces](assets/interfaces.png)

### ARP

| Metric Name | Type | Labels | Subsystem | Description | Disable Flag |
| --- | --- | --- | --- | --- | --- |
opnsense_arp_table_entries | Gauge | expired, hostname, interface_description, ip, mac, permanent, type | ARP Table | Arp entries by ip, mac, hostname, interface description, type, expired and permanent | --exporter.disable-arp-table |
opnsense_protocol_arp_sent_requests_total | Counter | n/a | Protocol Statistics | Total Number of sent ARP requests  | n/a |
opnsense_protocol_arp_received_requests_total | Counter | n/a | Protocol Statistics | Total Number of received ARP requests  | n/a |


### Gateways

| Metric Name | Type | Labels | Subsystem | Description | Disable Flag |
| --- | --- | --- | --- | --- | --- |
opnsense_gateways_status | Gauge | address, name | Gateways | Status of the gateway by name and address (1 = up, 0 = down, 2 = unknown) | n/a |
opnsense_gateways_loss_percentage | Gauge | address, name | Gateways | The current gateway loss percentage by name and address | n/a |
opnsense_gateways_rtt_milliseconds | Gauge | address, name | Gateways | RTT is the average (mean) of the round trip time in milliseconds by name and address | n/a |
opnsense_gateways_rttd_milliseconds | Gauge | address, name | Gateways | RTTd is the standard deviation of the round trip time in milliseconds by name and address | n/a |


### Protocol Statistics

| Metric Name | Type | Labels | Subsystem | Description | Disable Flag |
| --- | --- | --- | --- | --- | --- |
opnsense_protocol_tcp_sent_packets_total | Counter | n/a | Protocol Statistics | Total Number of sent TCP packets  | n/a |
opnsense_protocol_tcp_received_packets_total | Counter | n/a | Protocol Statistics | Total Number of received TCP packets | n/a |
opnsense_protocol_tcp_connection_count_by_state | Gauge | state | Protocol Statistics | Number of TCP connections by state | n/a |
opnsense_protocol_udp_delivered_packets_total | Counter | n/a | Protocol Statistics | Total Number of delivered UDP packets | n/a |
opnsense_protocol_udp_output_packets_total | Counter | n/a | Protocol Statistics | Total Number of output UDP packets  | n/a |
opnsense_protocol_udp_received_datagrams_total | Counter | n/a | Protocol Statistics | Total Number of received UDP Datagrams | n/a |
opnsense_protocol_udp_dropped_by_reason_total | CounterVector | reason | Protocol Statistics | Total Number of dropped UDP packets by reason | n/a |
opnsense_protocol_icmp_calls_total | Counter | n/a | Protocol Statistics | Total Number of ICMP calls | n/a |
opnsense_protocol_icmp_sent_packets_total | Counter | n/a | Protocol Statistics | Total Number of sent ICMP packets | n/a |
opnsense_protocol_icmp_dropped_by_reason_total | CounterVector | reason | Protocol Statistics | Total Number of dropped ICMP packets by reason | n/a |



### Unbound DNS

| Metric Name | Type | Labels | Subsystem | Description | Disable Flag |
| --- | --- | --- | --- | --- | --- |
opnsense_unbound_dns_uptime_seconds | Gauge | n/a | Unbound | Uptime of the unbound DNS service in seconds | --exporter.disable-unbound |

### Wireguard 

| Metric Name | Type | Labels | Subsystem | Description | Disable Flag |
| --- | --- | --- | --- | --- | --- |
opnsense_wireguard_interfaces_status | Gauge | name, description, public_key | Wireguard | Wireguard interfaces status by name, description and public key (1 = up, 0 = down) | --exporter.disable-wireguard  |
opnsense_wireguard_peer_received_bytes_total | Counter | device, device_type, device_name, peer_name | Wireguard | Bytes received by this wireguard peer | --exporter.disable-wireguard  |
opnsense_wireguard_peer_transmitted_bytes_total | Counter | device, device_type, device_name, peer_name | Wireguard | Bytes transmitted by this wireguard peer | --exporter.disable-wireguard  |
opnsense_wireguard_peer_last_handshake_seconds | Gauge | device, device_type, device_name, peer_name | Wireguard | Last handshake time in seconds by this wireguard peer | --exporter.disable-wireguard  |

### OpenVPN

| Metric Name | Type | Labels | Subsystem | Description | Disable Flag |
| --- | --- | --- | --- | --- | --- |
opnsense_openvpn_instances | Gauge | description, device_type, role, uuid | OpenVPN | OpenVPN instances (1 = enabled, 0 = disabled) by role (server, client) | --exporter.disable-openvpn |
