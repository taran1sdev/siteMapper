package main

import (
	"net/http"
	"net/url"
	"fmt"
	"flag"
	"strings"
	"link/link"
	"io"
)

var inputUrl = flag.String("url", "www.google.com", "The root of the domain to clone")

func main() {
	flag.Parse()
	
	links := get(*inputUrl)
	fmt.Println(links)
}

func get(link string)  []string {
	// Get the html response
	resp, err := http.Get(link)
	if err != nil{
		return []string{}
	}
	
	defer resp.Body.Close()
	
	// get the url from the response
	reqUrl := resp.Request.URL
	baseUrl := &url.URL {
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(getAnchors(resp.Body, base), withPrefix(base))
}

// Function to filter for local links - only links starting with the base URL will be returned

func filter(links []string, keepFn func(string) bool) (filtered []string) {
	for _, l := range links {
		if keepFn(l) {
			filtered = append(filtered, l)
		}
	}
	return
}

// Function returns a function that checks prefix of given string

func withPrefix(pfx string) func(string) bool {
	return func (link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

// Function Parses links from the html response and appends link to base domain if not present

func getAnchors(r io.Reader, base string) (anchors []string) {
	parsed, _ := link.ParseAnchors(r)

	for _, a := range parsed {
		switch {
			case strings.HasPrefix(a.Href, "/"):
				anchors = append(anchors, base+a.Href)
			case strings.HasPrefix(a.Href, "http"):
				anchors = append(anchors, a.Href)
		}
	}
	return
}
