package examples

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"time"
)

func RandomDelay() {
	url := "https://httpbin.org/delay/2"

	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "httpbin.*",
		RandomDelay: 5 * time.Second,
		Parallelism: 2,
	})

	for i := 0; i < 4; i++ {
		c.Visit(fmt.Sprintf("%s?n=%d", url, i))
	}

	c.Visit(url)
	c.Wait()
}
