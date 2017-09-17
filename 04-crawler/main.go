package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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

func crawl(url string, foundCh chan string) error {
	body, err := get(url)
	if err != nil {
		return err
	}
	for _, l := range parseLinks(body) {
		foundCh <- l
	}
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

	seeds = os.Args[1:]

	workCh := make(chan string)
	foundCh := make(chan string)
	visitedCh := make(chan string)
	doneCh := make(chan bool)

	fmt.Printf("Starting crawler with seeds: %s\n", strings.Join(seeds, ", "))

	// pushes seeds to the work channel
	go func() {
		for _, s := range seeds {
			workCh <- s
		}
	}()

	visited := make(map[string]bool)
	// receives found links and pushes them to the workCh (if they are new)
	// owns the "visited" map and knows when there is no work running anymore.
	// closes the doneCh if there is no work being performed or 10 links visited
	go func() {
		var currentWork int
		var totalVisited int
		for {
			select {
			case url := <-foundCh:
				if _, ok := visited[url]; ok {
					continue
				}
				visited[url] = false
				if totalVisited < 10 {
					currentWork++
					workCh <- url
				}
			case url := <-visitedCh:
				totalVisited++
				currentWork--
				visited[url] = true
				if currentWork == 0 {
					close(doneCh)
				}
			case <-doneCh:
				return
			}
		}
	}()

	// consumes workCh untill it receives on doneCh
	// posts links found to foundCh and visited urls to visitedCh
loop:
	for {
		select {
		case url := <-workCh:
			go func(url string) {
				crawl(url, foundCh)
				visitedCh <- url
			}(url)
		case <-doneCh:
			break loop
		}
	}

	// output all links found
	fmt.Println("Found links:")
	for k, v := range visited {
		fmt.Printf(" - %s (%s)\n", k, strconv.FormatBool(v))
	}

}
