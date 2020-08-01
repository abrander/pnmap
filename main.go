package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: os.Args[0],
	}

	interfaces *[]string

	hostInterfaces []net.Interface

	unknown string

	unknownFile   *os.File
	unknownWriter *pcapgo.Writer
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
		Use:     "monitor",
		Short:   "Monitor all interfaces for Probe Requests",
		Run:     monitor,
		PreRun:  setupWriter,
		PostRun: tearDownWriter,
	}
	monitorCmd.Flags().StringVarP(&unknown, "unknown", "u", "", "Path to write unknown packets to")
	rootCmd.AddCommand(monitorCmd)

	simulateCmd := &cobra.Command{
		Use:     "simulate",
		Short:   "",
		Run:     simulate,
		Args:    cobra.ExactArgs(1),
		PreRun:  setupWriter,
		PostRun: tearDownWriter,
	}
	simulateCmd.Flags().StringVarP(&unknown, "unknown", "u", "", "Path to write unknown packets to")
	rootCmd.AddCommand(simulateCmd)

	interfaces = monitorCmd.PersistentFlags().StringArrayP("interface", "i", []string{"all"}, "Interface(s) to monitor")
}

func setupWriter(_ *cobra.Command, _ []string) {
	var err error
	if unknown != "" {
		unknownFile, err = os.OpenFile(unknown, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf(err.Error())
		}

		unknownWriter = pcapgo.NewWriter(unknownFile)

		pos, _ := unknownFile.Seek(0, 2)
		if pos == 0 {
			err = unknownWriter.WriteFileHeader(65536, layers.LinkTypeEthernet)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
	}
}

func tearDownWriter(_ *cobra.Command, _ []string) {
	if unknownFile != nil {
		unknownFile.Close()
	}
}

func list(_ *cobra.Command, _ []string) {
	for _, i := range hostInterfaces {
		fmt.Printf("%s\n", i.Name)
	}

	os.Exit(0)
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

	i := newIntel()
	g := newGUI()

	i.hostChan = make(chan *NIC, 10)

	go func() {
		for packet := range packets {
			if !i.NewPacket(packet) && unknownWriter != nil {
				_ = unknownWriter.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
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

			filter(packet, packets)
		}
	}()

	go func() {
		for packet := range packets {
			fmt.Fprintf(os.Stderr, "%s\n", packet.String())
			if !i.NewPacket(packet) && unknownWriter != nil {
				_ = unknownWriter.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
			}
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

func filter(packet gopacket.Packet, out chan gopacket.Packet) {
	if ethernetLayer := packet.Layer(layers.LayerTypeEthernet); ethernetLayer != nil {
		eth := ethernetLayer.(*layers.Ethernet)
		ipv4 := packet.Layer(layers.LayerTypeIPv4)

		switch {
		// Throw away packets with no source.
		case eth.SrcMAC.String() == "00:00:00:00:00:00":

		// We're only interested in group traffic.
		case eth.DstMAC[0]&0x01 > 0:
			out <- packet

		// ... or IP broadcast traffic.
		case ipv4 != nil && ipv4.(*layers.IPv4).DstIP.Equal(net.IPv4bcast):
			out <- packet
		}
	}
}
