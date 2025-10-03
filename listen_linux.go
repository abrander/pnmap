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
		return
	}

	conn, err := raw.ListenPacket(intf, syscall.ETH_P_ALL, nil)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return
	}

	buffer := make([]byte, 65536)

	for {
		l, _, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Fatalf("error: %s", err.Error())
			break
		}

		packet := gopacket.NewPacket(buffer[0:l], layers.LayerTypeEthernet, gopacket.Default)
		packet.Metadata().Timestamp = time.Now()
		packet.Metadata().CaptureInfo.CaptureLength = l
		packet.Metadata().CaptureInfo.Length = l

		filter(packet, out)
	}
}
