pnmap
=====

Passive Network Mapper is an entirely passive network mapper. It will
passively and undetectable gather information about hosts and clients
participating in an ethernet segment.

Features
--------

- Undetectable by network participants
- Does not require promiscuous mode
- Supports wired and wireless networks
- Supports encrypted WiFi-networks
- Detects IPv4 addresses of hosts
- Detects IPv6 addresses of hosts
- Detects IPv6 neighbor discovery
- Detects IPv4 and IPv6 DHCP clients
- Detects public IPv4 address of natted network
- Detects DHCP hostnames
- Detects DHCP vendors
- Detects SSDP user agents
- Detects clients running Spotify and Spotify Connect speakers
- Detects Sonos speakers
- Detects Dropbox clients
- Detects HASP License Managers
- Detects MDNS services
- Detects macOS SSH servers
- Detects iOS and macOS hardware models
- Detects Chromecast and AirPlay clients and servers
- Detects various file-sharing services
- Detects Glen Dimplex Nobø Energy Control hubs
- Detects WS-Discovery clients
- Detects Ubiquiti Discover clients
- Detects TeamViewer
- Detects Minecraft clients
- Detects Steam
- Detects VNC
- Displays ethernet OUI vendors
- no cgo needed.

Requirements
------------

A working Go environment is required for compiling, and a Linux, BSD or
macOS host is required for running.

Compiling
---------

The usual `go mod download` and `go build` should suffice.

Running
-------
List network interfaces by invoking `./pnmap list`.

Monitoring a live network can be done like `./pnmap monitor -i eno1`.

Replaying a pcap file: `./pnmap simulate capture-file.pcap`.
