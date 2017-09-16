package main

import (
	"fmt"
	"os"
	"strings"

	_ "golang.org/x/net/html"
)

func main() {
	var seeds []string

	seeds = os.Args[1:]

	fmt.Printf("Starting crawler with seeds: %s\n", strings.Join(seeds, ", "))

	/* for each seed:
	get their content - net/http
	parse their content and find links golang.org/x/net/html
	store the links found
	*/

	// output all links found
}
