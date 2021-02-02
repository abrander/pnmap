package main

import (
	"encoding/binary"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var layerTypeMNDP = gopacket.RegisterLayerType(0x1001, gopacket.LayerTypeMetadata{Name: "MNDP", Decoder: gopacket.DecodeFunc(decodeMNDP)})

func init() {
	layers.RegisterUDPPortLayerType(layers.UDPPort(5678), layerTypeMNDP)
}

type MNDP struct {
	unknownHeader1 [2]byte
	unknownHeader2 [2]byte

	MAC        string
	Identity   string
	Version    string
	Platform   string
	Uptime     uint16
	SoftwareID string
	Board      string
	Interface  string

	contents []byte
}

func (m *MNDP) LayerType() gopacket.LayerType { return layerTypeMNDP }
func (m *MNDP) LayerContents() []byte         { return m.contents }
func (m *MNDP) LayerPayload() []byte          { return nil }

func decodeMNDP(data []byte, p gopacket.PacketBuilder) error {
	m := &MNDP{}

	// Well...
	copy(m.unknownHeader1[:], data)
	copy(m.unknownHeader2[:], data[2:])

	rest := data[4:]

	for len(rest) > 4 {
		var payload []byte
		typ := binary.BigEndian.Uint16(rest)
		length := binary.BigEndian.Uint16(rest[2:])

		rest = rest[4:]

		if len(rest) >= int(length) {
			payload = rest[:length]

			rest = rest[length:]
		}

		switch typ {
		case 1: // MAC address, string
			m.MAC = string(payload)
		case 5: // Identity, string
			m.Identity = string(payload)
		case 7: // Version, string
			m.Version = string(payload)
		case 8: // Platform, string
			m.Platform = string(payload)
		case 10: // Uptime, uint16
			m.Uptime = binary.BigEndian.Uint16(payload)
		case 11: // Software ID, string
			m.SoftwareID = string(payload)
		case 12: // Board, string
			m.Board = string(payload)
		case 13: // unknown, unseen
		case 14: // unknown, uint8
		case 15: // unknown, unseen
		case 16: // interface name, string
			m.Interface = string(payload)
		case 17: // unknown
		}
	}

	p.AddLayer(m)

	return p.NextDecoder(gopacket.DecodeFragment)
}
