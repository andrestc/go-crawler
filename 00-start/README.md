# A simple crawler

Let's start our first crawler implementation. This crawler will fetch a bunch of seed urls (provided 
as cli arguments) sequentially and output all the links found. A real crawler is much more complex but
this is a nice start.

## main.go

The file main.go, in this directory, has the initial code.

```go
func main() {
    var seeds []string

    seeds = os.Args[1:]

    fmt.Printf("Starting crawler with seeds: %s\n", strings.Join(seeds, ", "))
```

We begin by reading the seeds from the program arguments and print them to the user.

To try this out, simply run: `go run main.go https://andrestc.com`.

## Main loop

After reading the input, we should loop thru all `seeds` and print some message like "visiting https://adrestc.com".

Try implementing that now, using the `log` package instead of the `fmt` package.

You can check a package documentation on godoc: https://godoc.org/log

## Getting the content from a URL

Use the `net/http` package to make a request to the current URL and fetch it's content:

```go
response, err := http.Get("https://andrestc.com")
```

`http.Get` uses a global http client to make its requests and returns both the response and an error. Don't forget to check for errors properly.

## Parsing the HTML to find the links

Use the `golang.org/x/net/html` package to parse the content looking for links. The following sample explains how to check a single token to see if its a html tag (a link is a `a` tag):

```go
tk := html.NewTokenizer(response.Body)

token := tk.Next()
if token == html.StartTagToken {
    t := tk.Token()
    if t.Data != "a" {
        // not a link
    }
    for _, a := range t.Attr {
        if a.Key != "href" {
            continue
        }
        // a.Val is the link URL!
        fmt.Printf("found link: %s", a.Val)
        break
    }
}
```

Note: we should keep calling tk.Next() until the returned token is an `html.ErrorToken`, which means we finished parsing the document.

```go
if token == html.ErrorToken {
    // we finished parsing
}
```

Store all links found in a string slice.

## Output

Look thru the links found and print each one in a different line.