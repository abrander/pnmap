package main

import (
	"encoding/binary"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var layerTypeUDiscovery = gopacket.RegisterLayerType(0x1002, gopacket.LayerTypeMetadata{Name: "ubqt-discovery", Decoder: gopacket.DecodeFunc(decodeUbiquityDiscovery)})

func init() {
	layers.RegisterUDPPortLayerType(layers.UDPPort(10001), layerTypeUDiscovery)
}

type UDiscovery struct {
	Software string
	IP       string
	Uptime   int
	Name     string
	Model    string
	Series   string
}

func (m *UDiscovery) LayerType() gopacket.LayerType { return layerTypeUDiscovery }
func (m *UDiscovery) LayerContents() []byte         { return nil }
func (m *UDiscovery) LayerPayload() []byte          { return nil }

func decodeUbiquityDiscovery(data []byte, p gopacket.PacketBuilder) error {
	u := &UDiscovery{}

	dataLength := int(binary.BigEndian.Uint16(data[2:]))

	if dataLength+4 != len(data) {

		return p.NextDecoder(gopacket.DecodePayload)
	}

	rest := data[4:]

	for len(rest) >= 4 {
		var payload []byte

		typ := rest[0]
		length := binary.BigEndian.Uint16(rest[1:])

		rest = rest[3:]

		if len(rest) >= int(length) {
			payload = rest[:length]

			rest = rest[length:]
		} else {
			// malformed.
			return p.NextDecoder(gopacket.DecodePayload)
		}

		// Well, this is based on samples from a handful Ubiquiti
		// hardware. YMMV.
		switch typ {
		case 1: // MAC address, 6 bytes
		case 2: // MAC+IP? Legacy, 10 bytes
		case 3: // HW/SW, string
			u.Software = string(payload)
		case 4: // IP address, 4 bytes
			u.IP = net.IP(payload).String()
		case 10: // Uptime, uint32
			u.Uptime = int(binary.BigEndian.Uint32(payload))
		case 11: // name, string
			u.Name = string(payload)
		case 12: // model, string
			u.Model = string(payload)
		case 14: // 1 byte..?
		case 16: // 2 bytes..?
		case 20: // series, string
			u.Series = string(payload)
		}
	}

	p.AddLayer(u)

	return p.NextDecoder(gopacket.DecodeFragment)
}
