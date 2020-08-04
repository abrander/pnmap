package main

//go:generate go run contrib/embed-oui.go

var (
	privates = map[byte]bool{
		'2': true,
		'6': true,
		'a': true,
		'e': true,
	}
)

// OUIVendor will return the owner of a MAC address.
func OUIVendor(mac string) string {
	vendor := ouiToVendor[mac[0:8]]

	if vendor == "" && len(mac) > 1 && privates[mac[1]] {
		return "LAA (LOCALLY ADMINISTERED)"
	}

	return vendor
}
