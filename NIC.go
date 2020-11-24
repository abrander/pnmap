package main

import (
	"fmt"
	"time"
)

// NIC contains information about an ethernet station.
type NIC struct {
	MAC          string      `json:"MAC"`
	IPs          stringSlice `json:"IPs"`
	Hostnames    stringSlice `json:"Hostnames"`
	UserAgents   stringSlice `json:"UserAgents"`
	Vendor       stringSlice `json:"Vendor"`
	Applications stringSlice `json:"Applications"`
	Seen         int         `json:"Seen"`
	LastSeen     time.Time   `json:"LastSeen"`
	FirstSeen    time.Time   `json:"FirstSeen"`
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

	output += fmt.Sprintf("First seen: %s (%s ago)\n", n.FirstSeen.String(), time.Since(n.FirstSeen).String())
	output += fmt.Sprintf("Last seen: %s (%s ago)\n", n.LastSeen.String(), time.Since(n.LastSeen).String())
	output += fmt.Sprintf("Packets: %d\n\n", n.Seen)

	output += fmt.Sprintf("IPS: %v\n", n.IPs)
	output += fmt.Sprintf("Hostnames: %v\n", n.Hostnames)
	output += fmt.Sprintf("User agents: %v\n", n.UserAgents)
	output += fmt.Sprintf("Vendor: %v\n", n.Vendor)
	output += fmt.Sprintf("Applications: %v", n.Applications)

	return output
}
