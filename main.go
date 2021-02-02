package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

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

	unknown     string
	dissectOnly bool

	unknownFile   *os.File
	unknownWriter *pcapgo.Writer

	statefile = getStateFile()
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
		Args:    cobra.MinimumNArgs(1),
		PreRun:  setupWriter,
		PostRun: tearDownWriter,
	}
	simulateCmd.Flags().StringVarP(&unknown, "unknown", "u", "", "Path to write unknown packets to")
	simulateCmd.Flags().BoolVarP(&dissectOnly, "dissect-only", "d", false, "Only dissect packets")
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
		fmt.Printf("%10s %17s %s\n", i.Name, i.HardwareAddr.String(), OUIVendor(i.HardwareAddr.String()))
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

	state, _ := ioutil.ReadFile(statefile)
	json.Unmarshal(state, &i.NICCollection)

	i.hostChan = make(chan *NIC, 10)

	go func() {
		last := time.Now()

		for packet := range packets {
			if !i.NewPacket(packet) && unknownWriter != nil {
				_ = unknownWriter.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
			}
			now := time.Now()
			if now.Sub(last).Seconds() > 10 {
				os.Mkdir(filepath.Dir(statefile), 0700)
				f, _ := os.Create(statefile)
				j, _ := json.Marshal(i.NICCollection)
				_, _ = fmt.Fprintf(f, "%s", j)
				f.Close()
				last = time.Now()
			}
		}
	}()

	go func() {
		for nic := range i.hostChan {
			g.updateNIC(nic)
		}
	}()

	// Needs to be run in a Goroutine because of limited number of channels
	go func() {
		for _, nic := range i.NICCollection {
			i.hostChan <- nic
		}
		return
	}()

	_ = g.Run()
}

func simulate(_ *cobra.Command, args []string) {
	packets := make(chan gopacket.Packet, 10)

	i := newIntel()
	i.hostChan = make(chan *NIC, 10)

	go func() {
		for _, a := range args {
			f, err := os.Open(a)
			if err != nil {
				log.Fatalf(err.Error())
			}

			reader, err := pcapgo.NewReader(f)
			if err != nil {
				log.Fatalf(err.Error())
			}

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

			f.Close()
		}

		close(packets)
	}()

	if dissectOnly {
		go func() {
			for range i.hostChan {
			}
		}()

		for packet := range packets {
			if !i.NewPacket(packet) && unknownWriter != nil {
				_ = unknownWriter.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
			}
		}

		return
	}

	g := newGUI()

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

	err := g.Run()
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
		ipv6 := packet.Layer(layers.LayerTypeIPv6)

		switch {
		// Throw away packets with no source.
		case eth.SrcMAC.String() == "00:00:00:00:00:00":

		// We're only interested in group traffic.
		case eth.DstMAC[0]&0x01 > 0:
			out <- packet

		// ... or IPv4 broadcast traffic.
		case ipv4 != nil && ipv4.(*layers.IPv4).DstIP.Equal(net.IPv4bcast):
			out <- packet

		case ipv4 != nil && ipv4.(*layers.IPv4).DstIP.IsLinkLocalMulticast():
			out <- packet

		// ... or IPv6 broadcast traffic:
		case ipv6 != nil && ipv6.(*layers.IPv6).DstIP.IsLinkLocalMulticast():
			out <- packet
		}
	}
}
