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
- Detects clients running Spotify and Spotify Connect speakers
- Detects Dropbox client
- Detects HASP License Manager
- Detects MDNS services
- Detects iOS and macOS hardware model
- Detects Chromecast and AirPlay clients and servers
- Detects various file sharing services
- Displays ethernet OUI vendor
- no cgo needed.

Requirements
------------

A working Go environment is required for compiling and a Linux host is required for running.

Compiling
---------

The usual Go `go mod download` and `go build` should suffice.

Running
-------
List network interfaces by invoking `./pnmap list`.

Monitoring a live network can be done like `./pnmap monitor -i eno1`.

Replaying a pcap file: `./pnmap simulate capture-file.pcap`.
