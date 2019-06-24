package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Println("DO THE THING")
		os.Exit(0)
	}
	flag.Parse()
}