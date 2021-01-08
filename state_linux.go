package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getStateFile() string {
	homedir, _ := os.UserHomeDir()
	gateways := findGateways()

	if len(gateways) == 1 {
		mac := findMacFromIPInArpTable(gateways[0])
		return fmt.Sprintf("%s/.pnmap/state-%s-%s.json", homedir, mac, gateways[0])
	}

	// If we have multiple - or zero - gateways, we fall-back to a generic state file.
	return fmt.Sprintf("%s/.pnmap/state.json", homedir)
}

// find gateways
func findGateways() []string {
	file, err := os.Open("/proc/net/route")
	defer file.Close()

	if err != nil {
		log.Fatalf("Failed to open /proc/net/route")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var gateways []string
	scanner.Scan() // skip first line as it's a header
	for scanner.Scan() {
		s := strings.Fields(scanner.Text())
		if s[1] == "00000000" && s[7] == "00000000" {
			octet0, _ := strconv.ParseInt(s[2][6:8], 16, 64)
			octet1, _ := strconv.ParseInt(s[2][4:6], 16, 64)
			octet2, _ := strconv.ParseInt(s[2][2:4], 16, 64)
			octet3, _ := strconv.ParseInt(s[2][0:2], 16, 64)
			gateway := fmt.Sprintf("%d.%d.%d.%d", octet0, octet1, octet2, octet3)
			gateways = append(gateways, gateway)
		}
	}
	return gateways
}

// find the gateway in the arp table
func findMacFromIPInArpTable(ip string) string {
	file, err := os.Open("/proc/net/arp")
	defer file.Close()

	if err != nil {
		log.Fatalf("Failed to open /proc/net/route")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan() // skip first line as it's a header

	for scanner.Scan() {
		s := strings.Fields(scanner.Text())
		if s[0] == ip {
			return string(s[3]) // we don't expect to find more than one
		}
	}
	return ""
}
