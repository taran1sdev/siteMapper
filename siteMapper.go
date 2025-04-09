package main

import (
	"net/http"
	"net/url"
	"fmt"
	"flag"
	"strings"
	"link/link"
	"io"
	"encoding/xml"
	"os"
)

var inputUrl = flag.String("url", "https://www.google.com", "The root of the domain to map")
var maxDepth = flag.Int("depth", 2, "The maximum amount of links to follow")

// datatype for creating <loc></loc> tags in xml output

type loc struct {
	Loc string `xml:"loc"`  
}

// datatype for outer xml tags

type urlset struct {
	Urls []loc `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	flag.Parse()
	
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("  ", "    ")

	links := bfs(*inputUrl, *maxDepth)
	
	toXml := urlset{
		Xmlns:"http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	for _, link := range links {
		toXml.Urls = append(toXml.Urls, loc{Loc:link})
	}
	
	if err := enc.Encode(toXml); err != nil{
		fmt.Println("Error when encoding to xml", err)
	}
	fmt.Println()
}

// Function performs a breadth-first-search on all local links found in website

func bfs(urlString string, depth int) []string {
	seen := make(map[string]struct{})
	var queue map[string]struct{}
	nextQueue := map[string]struct{}{
		urlString: struct{}{},
	}
	for i := 0; i < depth; i++ {
		queue, nextQueue = nextQueue, make(map[string]struct{})
		if len(queue) == 0{
			break
		}
		for url, _ := range queue {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				nextQueue[link] = struct{}{}
			}
		}
	}
	links := make([]string, 0, len(seen))
	for url, _ := range seen{
		links = append(links, url)
	}
	return links
}

// Function performs a http GET requests to get html content to Parse

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

// Function returns a function that checks prefix of given string - For checking if links are local and start with http || /

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
