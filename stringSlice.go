package main

// stringSlice is a slice of strings containing unique strings.
type stringSlice []string

func (s *stringSlice) add(value string) {
	for _, v := range *s {
		if v == value {
			return
		}
	}

	*s = append(*s, value)
}
