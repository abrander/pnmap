package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/mdlayher/raw"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: os.Args[0],
	}

	interfaces *[]string

	hostInterfaces []net.Interface

	unknown string
)

func init() {
	hostInterfaces, _ = net.Interfaces()

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List network interfaces",
		Run:   list,
	}
	rootCmd.AddCommand(listCmd)

	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor all interfaces for Probe Requests",
		Run:   monitor,
	}
	monitorCmd.Flags().StringVarP(&unknown, "unknown", "u", "", "Path to write unknown packets to")
	rootCmd.AddCommand(monitorCmd)

	simulateCmd := &cobra.Command{
		Use:   "simulate",
		Short: "",
		Run:   simulate,
		Args:  cobra.ExactArgs(1),
	}
	rootCmd.AddCommand(simulateCmd)

	interfaces = monitorCmd.PersistentFlags().StringArrayP("interface", "i", []string{"all"}, "Interface(s) to monitor")
}

func list(_ *cobra.Command, _ []string) {
	for _, i := range hostInterfaces {
		fmt.Printf("%s\n", i.Name)
	}

	os.Exit(0)
}

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

func monitor(_ *cobra.Command, _ []string) {
	packets := make(chan gopacket.Packet, 10)

	if len(*interfaces) == 1 && (*interfaces)[0] == "all" {
		*interfaces = []string{}

		for _, i := range hostInterfaces {
			*interfaces = append(*interfaces, i.Name)
		}
	}

	for _, i := range *interfaces {
		go listen(i, packets)
	}

	var unknownWriter *pcapgo.Writer
	if unknown != "" {
		f, err := os.OpenFile(unknown, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf(err.Error())
		}
		defer f.Close()

		unknownWriter = pcapgo.NewWriter(f)

		pos, _ := f.Seek(0, 2)
		if pos == 0 {
			unknownWriter.WriteFileHeader(65536, layers.LinkTypeEthernet)
		}
	}

	i := newIntel()
	g := newGUI()

	i.hostChan = make(chan *NIC, 10)

	go func() {
		for packet := range packets {
			if !i.NewPacket(packet) && unknownWriter != nil {
				unknownWriter.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
			}
		}
	}()

	go func() {
		for nic := range i.hostChan {
			g.updateNIC(nic)
		}
	}()

	_ = g.Run()
}

func simulate(_ *cobra.Command, args []string) {
	packets := make(chan gopacket.Packet, 10)

	f, err := os.Open(args[0])
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer f.Close()

	reader, err := pcapgo.NewReader(f)
	if err != nil {
		log.Fatalf(err.Error())
	}

	i := newIntel()
	g := newGUI()

	i.hostChan = make(chan *NIC, 10)

	go func() {
		for {
			data, ci, err := reader.ReadPacketData()
			if err == io.EOF {
				break
			}

			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
			packet.Metadata().Timestamp = ci.Timestamp
			packet.Metadata().CaptureInfo = ci

			packets <- packet
		}
	}()

	go func() {
		for packet := range packets {
			fmt.Fprintf(os.Stderr, "%s\n", packet.String())
			i.NewPacket(packet)
		}
	}()

	go func() {
		for nic := range i.hostChan {
			g.updateNIC(nic)
		}
	}()

	err = g.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err.Error())
	}
}
