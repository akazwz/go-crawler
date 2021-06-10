package examples

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
)

func MaxDepth() {
	c := colly.NewCollector(
		colly.MaxDepth(1),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println(link)
		e.Request.Visit(link)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Fatal("something wrong with:", r.Request.URL, "\nError:", err)
	})

	c.Visit("https://github.com/")
}
