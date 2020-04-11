package main

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type mux map[gopacket.LayerType]func(source net.HardwareAddr, layer gopacket.Layer)

func newMux() mux {
	return make(mux)
}

func (m mux) add(layerType gopacket.LayerType, fun func(source net.HardwareAddr, layer gopacket.Layer)) {
	m[layerType] = fun
}

func (m mux) process(packet gopacket.Packet) {
	var source net.HardwareAddr

	if ethernetLayer := packet.Layer(layers.LayerTypeEthernet); ethernetLayer != nil {
		ethernet := ethernetLayer.(*layers.Ethernet)

		source = ethernet.SrcMAC
	}

	for t, f := range m {
		if l := packet.Layer(t); l != nil {
			f(source, l)
		}
	}
}
