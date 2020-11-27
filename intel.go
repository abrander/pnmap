package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/miekg/dns"
)

type intel struct {
	NICCollection map[string]*NIC
	hostChan      chan *NIC

	mux mux
}

func newIntel() *intel {
	i := &intel{
		NICCollection: make(map[string]*NIC),
		mux:           newMux(),
	}

	i.mux.add(layers.LayerTypeARP, i.arp)
	i.mux.add(layers.LayerTypeDHCPv4, i.dhcpv4)
	i.mux.add(layers.LayerTypeDHCPv6, i.dhcpv6)
	i.mux.add(layers.LayerTypeIPv4, i.ipv4)
	i.mux.add(layers.LayerTypeIPv6, i.ipv6)
	i.mux.add(layers.LayerTypeUDP, i.udp)
	i.mux.add(layers.LayerTypeICMPv6NeighborAdvertisement, i.ipv6NeighborAdvertisement)

	return i
}

func (i *intel) getNIC(addr []byte) *NIC {
	mac := mac(addr)

	nic, found := i.NICCollection[mac]
	if !found {
		n := newNIC(addr)
		i.NICCollection[mac] = n
		nic = n
	}

	i.hostChan <- nic

	return nic
}

func (i *intel) dhcpv4(source net.HardwareAddr, layer gopacket.Layer) bool {
	nic := i.getNIC(source)
	nic.Applications.add("dhcpv4")

	dhcpv4 := layer.(*layers.DHCPv4)
	if dhcpv4.Operation != layers.DHCPOpRequest {
		return false
	}

	for _, o := range dhcpv4.Options {

		switch o.Type {
		case layers.DHCPOptClassID:
			nic.Vendor.add(string(o.Data))
		case layers.DHCPOptHostname:
			nic.Hostnames.add(string(o.Data))
		}
	}

	return true
}

func (i *intel) dhcpv6(source net.HardwareAddr, layer gopacket.Layer) bool {
	nic := i.getNIC(source)
	nic.Applications.add("dhcpv6")

	return true
}

func (i *intel) arp(source net.HardwareAddr, layer gopacket.Layer) bool {
	arp := layer.(*layers.ARP)

	nic := i.getNIC(source)

	if len(arp.SourceProtAddress) == 4 {
		ip := fmt.Sprintf("%d.%d.%d.%d", arp.SourceProtAddress[0], arp.SourceProtAddress[1], arp.SourceProtAddress[2], arp.SourceProtAddress[3])
		if ip != "0.0.0.0" {
			nic.IPs.add(ip)
		}

		return true
	}

	return false
}

func (i *intel) ipv6(source net.HardwareAddr, layer gopacket.Layer) bool {
	ipv6 := layer.(*layers.IPv6)

	nic := i.getNIC(source)

	if ip := ipv6.SrcIP.String(); ip != "::" {
		nic.IPs.add(ip)
	}

	return false
}

func (i *intel) ipv4(source net.HardwareAddr, layer gopacket.Layer) bool {
	ipv4 := layer.(*layers.IPv4)

	nic := i.getNIC(source)

	ip := ipv4.SrcIP.String()
	if ip != "0.0.0.0" {
		nic.IPs.add(ip)
	}

	return false
}

func (i *intel) udp(source net.HardwareAddr, layer gopacket.Layer) bool {
	udp := layer.(*layers.UDP)
	nic := i.getNIC(source)

	switch udp.DstPort {

	// NBNS
	case 137:
		nic.Applications.add("NetBIOS-Name-Service")
		return true

	// NBDS - SMB
	case 138:
		nic.Applications.add("NetBIOS-Datagram-Service")
		return true

	// SSDP
	case 1900:
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(udp.Payload)))
		if err != nil {
			return false
		}

		ua := req.Header.Get("user-agent")
		if ua != "" {
			nic.UserAgents.add(ua)
		}

		return true

	// HASP License Manager
	case 1947:
		nic.Applications.add("HASP-License-Manager")

		return true

	// WS-Discovery
	case 3702:
		nic.Applications.add("WS-Discovery")

		return true

	// Multicast-DNS
	case 5353:
		msg := new(dns.Msg)

		dnsParts := func(in string) []string {
			in = strings.TrimSuffix(in, ".local.")
			parts := []string{""}
			part := 0

			var r rune
			for i, w := 0, 0; i < len(in); i += w {
				r, w = utf8.DecodeRuneInString(in[i:])
				if r == '\\' {
					var w2 int
					r, w2 = utf8.DecodeRuneInString(in[i+w:])
					w += w2
					parts[part] += string(r)
				} else if r == '.' {
					parts = append(parts, "")
					part++
				} else {
					parts[part] += string(r)
				}
			}

			return parts
		}

		if err := msg.Unpack(udp.Payload); err != nil {
			return false
		}

		if !msg.Response {
			return true
		}

		m := map[string]string{
			"_sftp-ssh":        "SSH",
			"_smb":             "Samba",
			"_ipp":             "IPP",
			"_ipps":            "IPPS",
			"_pdl-datastream":  "PDL-socket",
			"_afpovertcp":      "AFP",          // Apple Filing Protocol
			"_raop":            "AirPlay-RAOP", // Remote Audio Output Protocol
			"_airplay":         "AirPlay-display",
			"_companion-link":  "AirPlay-client",
			"_services":        "",
			"_nvstream_dbd":    "NVidia-Gamestream",
			"_homekit":         "homekit?",
			"_ePCL":            "ePCL?",
			"_universal":       "universal?",
			"_print":           "print?",
			"_wfds-print":      "wfds-print?",
			"_printer":         "LPR-printer",
			"_http":            "HTTP-server",
			"_scanner":         "Scanner",
			"_http-alt":        "HTTP-server-alt",
			"_uscan":           "uscan?",
			"_privet":          "Privet",
			"_uscans":          "uscans?",
			"_soundtouch":      "SoundTouch", // Bose
			"_googlecast":      "Chromecast",
			"_spotify-connect": "Spotify-Connect",
			"_teamviewer":      "TeamViewer",
			"_rfb":             "VNC",
			"_adisk":           "TimeCapsule",
			"_telnet":          "Telnet",
			"_sonos":           "Sonos",
		}

		for _, answer := range msg.Answer {
			names := dnsParts(answer.Header().Name)
			switch rr := answer.(type) {
			case *dns.A:
				name := strings.TrimSuffix(rr.Header().Name, ".local.")

				nic.Hostnames.add(name)

			case *dns.PTR:
				app, found := m[names[0]]
				if !found {
					app = names[0]
				}

				if strings.HasSuffix(rr.Header().Name, ".arpa.") {
					break
				}

				if app != "" {
					nic.Applications.add(app)
				}

			case *dns.SRV:
				if len(names) < 2 {
					break
				}

				app, found := m[names[1]]
				if !found {
					app = names[1]
				}

				if app != "" {
					nic.Applications.add(app)
				}

				if names[0][0] != '_' {
					nic.Hostnames.add(names[0])
				}

			case *dns.TXT:
				nic.Hostnames.add(names[0])
				if names[1] == "_device-info" {
					nic.Vendor.add(rr.Txt[0])
				}
			}
		}

	case 10000, 10001:
		// Nobø Hub
		// https://www.glendimplex.se/media/15650/nobo-hub-api-v-1-1-integration-for-advanced-users.pdf
		if bytes.Contains(udp.Payload, []byte("__NOBOHUB__")) {
			nic.Vendor.add("Glen-Dimplex")
			nic.Applications.add("nobo")

			return true
		}

		// Ubiquiti discover clients
		l := len(udp.Payload)
		if udp.DstPort == 10001 && l > 3 {
			if plen := udp.Payload[3]; int(plen)+4 == l {
				nic.Applications.add("ubnt-discover")

				return true
			}
		}

	// Dropbox
	case 17500:
		dummy := make(map[string]interface{})
		err := json.Unmarshal(udp.Payload, &dummy)
		if err == nil {
			// If we can decode a JSON payload, we assume it's
			// from Dropbox.
			nic.Applications.add("Dropbox")

			return true
		}

	// Raknet for Minecraft client
	case 19133:
		nic.Applications.add("Minecraft")

		return true

	// Steam client
	case 27036:
		nic.Applications.add("Steam")

		if len(udp.Payload) < 40 {
			return false
		}

		hlen, err := binary.ReadUvarint(bytes.NewBuffer(udp.Payload[8:]))
		if err != nil {
			return false
		}

		// 8: skip 8 byts of signature
		buf := proto.NewBuffer(udp.Payload[8:])
		if err != nil {
			return false
		}

		// signature + header length + header + body length
		offset := 8 + 4 + hlen + 4

		if offset > uint64(len(udp.Payload)) {
			return false
		}

		buf.SetBuf(udp.Payload[offset:])

		for value, err := buf.DecodeVarint(); err == nil; value, err = buf.DecodeVarint() {
			number := value >> 3
			typ := value & 0x7

			var str string

			switch typ {
			case 0:
				_, err = buf.DecodeVarint()
			case 1:
				_, err = buf.DecodeFixed64()
			case 2:
				str, err = buf.DecodeStringBytes()
			case 5:
				_, err = buf.DecodeFixed32()
			default:
				return false
			}

			if err != nil {
				break
			}

			switch number {
			case 4:
				nic.Hostnames.add(str)

			case 20, 21:
				if str != "0.0.0.0" {
					nic.IPs.add(str)
				}
			}
		}

		return true

	// Spotify
	case 57621:
		if bytes.HasPrefix(udp.Payload, []byte("SpotUdp")) {
			nic.Applications.add("Spotify")

			return true
		}
	}

	return false
}

func (i *intel) ipv6NeighborAdvertisement(source net.HardwareAddr, layer gopacket.Layer) bool {
	na := layer.(*layers.ICMPv6NeighborAdvertisement)
	nic := i.getNIC(source)

	nic.IPs.add(na.TargetAddress.String())

	return true
}

func (i *intel) NewPacket(packet gopacket.Packet) bool {
	if ethernetLayer := packet.Layer(layers.LayerTypeEthernet); ethernetLayer != nil {
		ethernet := ethernetLayer.(*layers.Ethernet)

		nic := i.getNIC(ethernet.SrcMAC)

		nic.LastSeen = packet.Metadata().Timestamp
		if nic.FirstSeen.IsZero() {
			nic.FirstSeen = packet.Metadata().Timestamp
		}
		nic.Seen++

		return i.mux.process(packet)
	}

	return false
}
