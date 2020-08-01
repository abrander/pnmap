// +build darwin dragonfly freebsd netbsd openbsd

package main

import (
	"log"
	"syscall"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/bsdbpf"
	"github.com/google/gopacket/layers"
)

func listen(deviceName string, out chan gopacket.Packet) {
	options := &bsdbpf.Options{
		ReadBufLen:       32767,
		Promisc:          false,
		Immediate:        true,
		PreserveLinkAddr: true,
	}

	sniffer, err := bsdbpf.NewBPFSniffer(deviceName, options)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	for {
		buffer, ci, err := sniffer.ReadPacketData()
		if err != nil {
			if e, ok := err.(syscall.Errno); ok && e.Temporary() {
				continue
			}

			// This will happen from time to time - we simply continue and hope for the best :-)
			if err.Error() == "BPF captured frame received with corrupted BpfHdr struct." {
				continue
			}

			log.Fatalf("ReadPacketData: %s", err.Error())
		}

		packet := gopacket.NewPacket(buffer[0:ci.CaptureLength], layers.LayerTypeEthernet, gopacket.Default)
		packet.Metadata().CaptureInfo = ci
		packet.Metadata().Timestamp = time.Now()

		// If we're unable to decode the packet, continue in silence.
		_, failure := packet.Layer(gopacket.LayerTypeDecodeFailure).(*gopacket.DecodeFailure)
		if failure {
			continue
		}

		filter(packet, out)
	}
}
