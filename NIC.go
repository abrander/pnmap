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
	output := ""

	output += fmt.Sprintf("[yellow]First seen[reset]: [white]%s[reset] ([white]%s[reset] ago)\n", n.FirstSeen.UTC().String(), time.Since(n.FirstSeen).Round(time.Second).String())
	output += fmt.Sprintf("[yellow]Last seen[reset]: [white]%s[reset] ([white]%s[reset] ago)\n", n.LastSeen.UTC().String(), time.Since(n.LastSeen).Round(time.Second).String())
	output += fmt.Sprintf("[yellow]Packets[reset]: [white]%d[reset]\n\n", n.Seen)

	output += fmt.Sprintf("[yellow]OUI Vendor[reset]: [white]%s[reset]\n", OUIVendor(n.MAC))
	output += fmt.Sprintf("[yellow]IPS[reset][reset]: %s\n", n.IPs.String())
	output += fmt.Sprintf("[yellow]Hostnames[reset]: %s\n", n.Hostnames.String())
	output += fmt.Sprintf("[yellow]User agents[reset]: %s\n", n.UserAgents.String())
	output += fmt.Sprintf("[yellow]Vendor[reset]: %s\n", n.Vendor.String())
	output += fmt.Sprintf("[yellow]Applications[reset]: %s", n.Applications.String())

	return output
}
