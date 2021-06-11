package examples

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

func Queue() {
	url := "https://httpbin.org/delay/1"

	c := colly.NewCollector()

	q, _ := queue.New(
		2,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	for i := 0; i < 5; i++ {
		q.AddURL(fmt.Sprintf("%s?n=%d", url, i))
	}

	q.Run(c)
}
