package main

import (
	"log"
	"net"
	"syscall"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/mdlayher/raw"
)

func listen(deviceName string, out chan gopacket.Packet) {
	intf, err := net.InterfaceByName(deviceName)
	if err != nil {
		log.Printf("error: %s", err.Error())
	}

	conn, err := raw.ListenPacket(intf, syscall.ETH_P_ALL, nil)
	if err != nil {
		log.Printf("error: %s", err.Error())
	}

	buffer := make([]byte, 65536)

	for {
		l, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Fatalf("error: %s", err.Error())
			break
		}

		// Throw away packets with no source.
		if addr.String() == "00:00:00:00:00:00" {
			continue
		}

		packet := gopacket.NewPacket(buffer[0:l], layers.LayerTypeEthernet, gopacket.Default)
		packet.Metadata().Timestamp = time.Now()
		packet.Metadata().CaptureInfo.CaptureLength = l
		packet.Metadata().CaptureInfo.Length = l

		if ethernetLayer := packet.Layer(layers.LayerTypeEthernet); ethernetLayer != nil {
			eth := ethernetLayer.(*layers.Ethernet)

			// We're only interested in group traffic.
			if eth.DstMAC[0]&0x01 > 0 {
				out <- packet
			}
		}
	}
}
