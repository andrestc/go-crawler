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

func crawl(url string, foundCh chan string, done chan bool) error {
	body, err := get(url)
	if err != nil {
		return err
	}
	for _, l := range parseLinks(body) {
		foundCh <- l
	}
	done <- true
	return body.Close()
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

	foundCh := make(chan string)
	doneCh := make(chan bool)
	fmt.Printf("Starting crawler with seeds: %s\n", strings.Join(seeds, ", "))

	for _, s := range seeds {
		go crawl(s, foundCh, doneCh)
	}

	var totalDone int
	for totalDone < len(seeds) {
		select {
		case url := <-foundCh:
			links = append(links, url)
		case <-doneCh:
			totalDone++
		}
	}

	// output all links found
	fmt.Println("Found links:")
	for _, l := range links {
		fmt.Printf(" - %s\n", l)
	}

}
