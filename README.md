pnmap
=====

Passive Network Mapper is a completely passive network mapper. It will
passively and undetectable gather information about hosts participating in an ethernet network
segment.

Features
--------

- Undetectable by network participants
- Does not require promiscuous mode
- Detects IPv4 addresses of hosts
- Detects IPv6 addresses of hosts
- Detects DHCP hostnames
- Detects DHCP vendors
- Detects SSDP user agents
- Displays ethernet OUI vendor

Requirements
------------

A working Go environment is required for compiling.

pnmap uses libpcap and requires the `libpcap0.8-dev` package.
