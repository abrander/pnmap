package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type ouiDatabase struct {
	toVendor map[string]string
}

var (
	privates = map[byte]bool{
		'2': true,
		'6': true,
		'a': true,
		'e': true,
	}
)

func newOuiDatabase() (*ouiDatabase, error) {
	o := &ouiDatabase{
		toVendor: make(map[string]string),
	}

	f, err := os.Open("oui.csv")
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)

	// Burn the header.
	_, _ = r.Read()

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if len(record[1]) != 6 {
			panic(record)
		}

		lower := strings.ToLower(record[1])
		mac := fmt.Sprintf("%c%c:%c%c:%c%c", lower[0], lower[1], lower[2], lower[3], lower[4], lower[5])
		vendor := strings.TrimSpace(record[2])

		o.toVendor[mac] = vendor
	}

	return o, nil
}

func (o *ouiDatabase) Vendor(mac string) string {
	vendor := o.toVendor[mac[0:8]]

	if vendor == "" && len(mac) > 1 && privates[mac[1]] {
		return "LAA (LOCALLY ADMINISTERED)"
	}

	return vendor
}
