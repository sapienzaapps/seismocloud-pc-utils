// +build nogui

package main

import (
	"flag"
	"fmt"
	"os"
)

func guimain() {
	fmt.Printf("This version of %s was built without GUI\n\n", os.Args[0])
	flag.PrintDefaults()
}
