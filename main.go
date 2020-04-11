package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: os.Args[0],
	}

	interfaces *[]string

	hostInterfaces []net.Interface
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
	handle, err := pcap.OpenLive(deviceName, 65535, false, pcap.BlockForever)
	if err != nil {
		fmt.Printf("Error opening device %s: %s\n", deviceName, err.Error())
		os.Exit(1)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for {
		packet, err := packetSource.NextPacket()
		if err != nil {
			fmt.Printf("Error on %s: %s\n", deviceName, err.Error())
			os.Exit(1)
		} else if err == nil {
			out <- packet
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

	i := newIntel()
	g := newGUI()

	i.hostChan = make(chan *NIC, 10)

	go func() {
		for {
			select {
			case packet := <-packets:
				go i.NewPacket(packet)
			case nic := <-i.hostChan:
				g.updateNIC(nic)
			}
		}
	}()

	_ = g.Run()
}

func simulate(_ *cobra.Command, args []string) {
	packets := make(chan gopacket.Packet, 10)

	handle, err := pcap.OpenOffline(args[0])
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	i := newIntel()
	g := newGUI()

	i.hostChan = make(chan *NIC, 10)

	go func() {
		for {
			packet, err := packetSource.NextPacket()
			if err != nil && err == io.EOF {
				break
			}

			packets <- packet
		}
	}()

	go func() {
		for {
			select {
			case packet := <-packets:
				go i.NewPacket(packet)
			case nic := <-i.hostChan:
				g.updateNIC(nic)
			}
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
