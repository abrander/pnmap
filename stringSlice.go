package main

import (
	"strings"
)

// stringSlice is a slice of strings containing unique strings.
type stringSlice []string

func (s *stringSlice) add(value string) {
	if len(value) == 0 {
		return
	}

	for _, v := range *s {
		if v == value {
			return
		}
	}

	*s = append(*s, value)
}

func (s *stringSlice) String() string {
	all := strings.Join([]string(*s), "[reset], [white]")

	return "[reset][[white]" + all + "[reset]]"
}
