package examples

import (
	"fmt"
	"github.com/gocolly/colly"
)

func Parallel() {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		fmt.Println(link)

		e.Request.Visit(link)
	})

	c.Visit("https://github.com/")
	c.Wait()
}
