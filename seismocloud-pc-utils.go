package main

import (
	"flag"
)

const (
	PKTTYPE_DISCOVERY       = 1
	PKTTYPE_DISCOVERY_REPLY = 2
)

func main() {
	discoverFlag := flag.Bool("discover", false, "Perform discovery and exit")

	flag.Parse()

	if *discoverFlag {
		discovery(0)
	} else {
		guimain()
	}
}
