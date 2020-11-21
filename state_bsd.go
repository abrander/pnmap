// +build darwin dragonfly freebsd netbsd openbsd

package main

import (
	"fmt"
	"os"
)

func getStateFile() string {
	homedir, _ := os.UserHomeDir()
	return fmt.Sprintf("%s/.pnmap/state.json", homedir)
}
