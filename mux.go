package main

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type mux map[gopacket.LayerType]func(source net.HardwareAddr, layer gopacket.Layer) bool

func newMux() mux {
	return make(mux)
}

func (m mux) add(layerType gopacket.LayerType, fun func(source net.HardwareAddr, layer gopacket.Layer) bool) {
	m[layerType] = fun
}

func (m mux) process(packet gopacket.Packet) bool {
	var source net.HardwareAddr

	if ethernetLayer := packet.Layer(layers.LayerTypeEthernet); ethernetLayer != nil {
		ethernet := ethernetLayer.(*layers.Ethernet)

		source = ethernet.SrcMAC
	}

	recognized := false
	for t, f := range m {
		if l := packet.Layer(t); l != nil {
			r := f(source, l)
			if r {
				recognized = true
			}
		}
	}

	return recognized
}
