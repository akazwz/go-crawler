package examples

import (
	"bytes"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"log"
)

func ProxySwitcher() {
	c := colly.NewCollector(colly.AllowURLRevisit())

	rp, err := proxy.RoundRobinProxySwitcher(
		"socks5://127.0.0.1:1337",
		"socks5://127.0.0.1:1338",
	)
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

	c.OnResponse(func(r *colly.Response) {
		log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
	})

	for i := 0; i < 5; i++ {
		err := c.Visit("https://httpbin.org/ip")
		if err != nil {
			log.Fatal(err)
		}
	}
}
