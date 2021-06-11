package examples

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
)

func RequestContext() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.Ctx.Get("url"))
	})

	err := c.Visit("https://baidu.com/")
	if err != nil {
		log.Fatal(err)
	}
}
