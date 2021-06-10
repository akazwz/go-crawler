package examples

import (
	"fmt"
	"github.com/gocolly/colly"
)

func ErrorHandling() {
	c := colly.NewCollector()

	c.OnHTML("*", func(e *colly.HTMLElement) {
		fmt.Println(e)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL,
			"\nfailed with response:", r,
			"\nError:", err)
	})

	c.Visit("https://definitely-not-a.website/")
}
