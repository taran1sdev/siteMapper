# siteMapper

This application uses the link module in /htmlLinkParser. It takes a url as input, performs a GET request on the url, extracts all the `<a>` tags from html response and filters those to only include links within the same domain.

The application then performs a breadth-first-search of the local links found on the domain - the user can specify the depth of this search.

The application then encodes the links found into an xml format following the Sitemap protocol schema found at :- https://www.sitemaps.org/protocol.html.

The application will write the xml to a file if one is provided, otherwise it will write to Stdout.

Build:
```bash
git clone https://github.com/taran1sdev/siteMapper.git
go mod init link
go build siteMapper.go
```

Usage:
```bash
./siteMapper -url=https://some-url.com -depth=3 -o=out.xml
```

Example:
```bash
$ go run siteMapper.go -url=https://go.dev/doc/ -depth=2  
  <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
      <url>
          <loc>https://go.dev/doc/codewalk/markov</loc>
      </url>
      <url>
          <loc>https://go.dev/doc/tutorial/workspaces.html</loc>
      </url>
      <url>
          <loc>https://go.dev/blog/v2-go-modules</loc>
      </url>
      <url>
          <loc>https://go.dev/cmd/go/</loc>
      </url>
      <url>
          <loc>https://go.dev/ref/mem</loc>
      </url>
...SNIP
      <url>
          <loc>https://go.dev/doc</loc>
      </url>
      <url>
          <loc>https://go.dev/doc/devel/release</loc>
      </url>
  </urlset>
```
