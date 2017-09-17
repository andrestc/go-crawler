package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func get(url string) (io.ReadCloser, error) {
	log.Printf("Visiting %s", url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func parseLinks(r io.Reader) []string {
	tk := html.NewTokenizer(r)
	var links []string
	for {
		tt := tk.Next()

		switch tt {
		case html.ErrorToken:
			return links
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

func main() {
	var seeds []string
	var links []string

	seeds = os.Args[1:]

	fmt.Printf("Starting crawler with seeds: %s\n", strings.Join(seeds, ", "))

	// for each seed:
	for _, s := range seeds {
		body, err := get(s)
		if err != nil {
			log.Printf("Error crawling %s: %v", s, err)
			continue
		}
		links = append(links, parseLinks(body)...)
		err = body.Close()
		if err != nil {
			log.Printf("Failed to close body: %v", err)
		}
	}

	// output all links found
	fmt.Println("Found links:")
	for _, l := range links {
		fmt.Printf(" - %s\n", l)
	}

}
