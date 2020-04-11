package main

import (
	"fmt"
	"time"
)

// NIC contains information about an ethernet station.
type NIC struct {
	MAC        string
	IPs        stringSlice
	Hostnames  stringSlice
	userAgents stringSlice
	vendor     stringSlice
	seen       int
	lastSeen   time.Time
}

func mac(addr []byte) string {
	if len(addr) == 0 {
		return ""
	}

	mac := fmt.Sprintf("%02x", addr[0])
	for i := 1; i < len(addr); i++ {
		mac += fmt.Sprintf(":%02x", addr[i])
	}

	return mac
}

func newNIC(addr []byte) *NIC {
	return &NIC{MAC: mac(addr)}
}

func (n *NIC) String() string {
	output := n.MAC + "\n\n"

	output += fmt.Sprintf("Last seen: %s (%s ago)\n\n", n.lastSeen.String(), time.Since(n.lastSeen).String())

	output += fmt.Sprintf("IPS: %v\n", n.IPs)
	output += fmt.Sprintf("Hostnames: %v\n", n.Hostnames)
	output += fmt.Sprintf("User agents: %v\n", n.userAgents)
	output += fmt.Sprintf("Vendor: %v\n", n.vendor)

	return output
}
