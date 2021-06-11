package examples

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

func RateLimit() {
	url := "https://httpbin.org/delay/2"

	c := colly.NewCollector(
		colly.Async(true),
		colly.Debugger(&debug.LogDebugger{}),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 2,
	})

	for i := 0; i < 5; i++ {
		c.Visit(fmt.Sprintf("%s?n=%d", url, i))
	}

	c.Wait()
}
