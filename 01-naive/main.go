package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// This is our first solution.
// Code works but has 2 problems:
// - Not concurrent, we can do better if we fetch pages in a concurrent fashion
// - Code is all in main, we should refactor
func main() {
	var seeds []string
	var links []string

	seeds = os.Args[1:]

	fmt.Printf("Starting crawler with seeds: %s\n", strings.Join(seeds, ", "))

	// for each seed:
	for _, s := range seeds {
		log.Printf("Visiting %s", s)
		//get their content - net/http
		response, err := http.Get(s)
		if err != nil {
			log.Printf("Failed to get %s: %v", s, err)
			continue
		}

		//parse their content and find links golang.org/x/net/html
		tk := html.NewTokenizer(response.Body)

	loop:
		for {
			tt := tk.Next()

			switch tt {
			case html.ErrorToken:
				response.Body.Close()
				break loop
			case html.StartTagToken:
				t := tk.Token()
				if t.Data != "a" {
					break
				}
				for _, a := range t.Attr {
					if a.Key != "href" {
						continue
					}
					// store the links found
					links = append(links, a.Val)
					break
				}
			}
		}
	}

	// output all links found
	fmt.Println("Found links:")
	for _, l := range links {
		fmt.Printf(" - %s\n", l)
	}

}
