package main

import (
	"fmt"
	"strings"
)

func main() {
	var seeds []string

	// read seeds from command line arguments

	fmt.Println("Starting crawler with seeds: %s", strings.Join(seeds, ", "))

	// start go routines to follow seeds

	// output all pages visited
}
