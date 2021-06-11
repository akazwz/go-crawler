package examples

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"regexp"
)

func URLFilter() {
	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile("https://httpbin\\.org/(|e.+)$"),
			regexp.MustCompile("https://httpbin\\.org/h.+"),
		),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := c.Visit("https://httpbin.org")
	if err != nil {
		log.Fatal(err)
	}

}
