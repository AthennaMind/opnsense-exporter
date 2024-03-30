## OPNsense Exporter Metrics List

This table represent each metric and it's labels, the subsystem that it belongs, its description and how to disable it.


| Metric Name | Type | Labels | Subsystem | Description | Disable Flag |
| --- | --- | --- | --- | --- | --- |
opnsense_exporter_scrapes_total | Counter | n/a | n/a | Total number of scrapes by the OPNsense exporter | n/a |
opnsense_exporter_endpoint_errors_total | Counter | endpoint | n/a | Total number of errors by endpoint returned by the OPNsense API during data fetching | n/a |
opnsense_arp_table_entries | Gauge | expired, hostname, interface_description, ip, mac, permanent, type | ARP Table | Arp entries by ip, mac, hostname, interface description, type, expired and permanent | --exporter.disable-arp-table |
opnsense_cron_job_status | Gauge | command, description, origin, schedule | Cron Table | Cron job status by name and description (1 = enabled, 0 = disabled) | --exporter.disable-cron-table |
opnsense_gateways_status | Gauge | address, name | Gateways | Status of the gateway by name and address (1 = up, 0 = down, 2 = unknown) | n/a |
opnsense_gateways_loss_percentage | Gauge | address, name | Gateways | The current gateway loss percentage by name and address | n/a |
opnsense_gateways_rtt_milliseconds | Gauge | address, name | Gateways | RTT is the average (mean) of the round trip time in milliseconds by name and address | n/a |
opnsense_gateways_rttd_milliseconds | Gauge | address, name | Gateways | RTTd is the standard deviation of the round trip time in milliseconds by name and address | n/a |
opnsense_openvpn_instances | Gauge | description, device_type, role, uuid | OpenVPN | OpenVPN instances (1 = enabled, 0 = disabled) by role (server, client) | n/a |
opnsense_protocol_arp_sent_requests_total | Counter | n/a | Protocol Statistics | Total Number of sent ARP requests by the system | n/a |
opnsense_protocol_arp_received_requests_total | Counter | n/a | Protocol Statistics | Total Number of received ARP requests by the system | n/a |
opnsense_protocol_tcp_sent_packets_total | Counter | n/a | Protocol Statistics | Total Number of sent TCP packets by the system | n/a |
opnsense_protocol_tcp_received_packets_total | Counter | n/a | Protocol Statistics | Total Number of received TCP packets by the system | n/a |
opnsense_protocol_tcp_connection_count_by_state | Gauge | state | Protocol Statistics | Number of TCP connections by state | n/a |
opnsense_services_running_total | Gauge | n/a | Services | Total number of running services | n/a |
opnsense_services_stopped_total | Gauge | n/a | Services | Total number of stopped services | n/a |
opnsense_services_status | Gauge | name, description | Services | Service status by name and description (1 = running, 0 = stopped) | n/a |
opnsense_unbound_dns_uptime_seconds | Gauge | n/a | Unbound | Uptime of the unbound DNS service in seconds | --exporter.disable-unbound |