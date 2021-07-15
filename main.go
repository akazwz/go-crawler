package main

import (
	"fmt"
	"github.com/akazwz/go-crawler/douban"
)

func main() {
	fmt.Println("Hello, Colly!")
	// examples.Basic()
	// examples.ErrorHandling()
	// examples.Login()
	// examples.MaxDepth()
	// examples.Multipart()
	// examples.Parallel()
	// examples.ProxySwitcher()
	// examples.Queue()
	// examples.RandomDelay()
	// examples.RateLimit()
	// examples.RedisBackend()
	// examples.RequestContext()
	// examples.ScraperServer()
	// examples.URLFilter()
	// real_life_examples.CryptocoinsMarketCapacity()
	douban.Movies()
	// douban.Books()
	// weibo.HotSearch()
}
