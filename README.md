pnmap
=====

Passive Network Mapper is an entirely passive network mapper. It will
passively and undetectable gather information about hosts participating
in an ethernet network segment.

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

Compiling
---------

The usual Go `go get -u ./...` and `go buiuld .` should suffice.

Running
-------
List network interfaces by invoking `./pnmap list`.

Monitoring a live network can be done like `./pnmap monitor -i eno1`.

Replaying a pcap file: `./pnmap simulate capture-file.pcap`.
