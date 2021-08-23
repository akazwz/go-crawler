package main

import (
	"fmt"
	"github.com/akazwz/go-crawler/global"
	"github.com/akazwz/go-crawler/initialize"
	"github.com/akazwz/go-crawler/weibo"
	"github.com/jasonlvhit/gocron"
	"log"
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
	// douban.Movies()
	// douban.Books()

	global.VP = initialize.InitViper()
	if global.VP == nil {
		fmt.Println("配置文件初始化失败")
	}
	fmt.Println(global.CFG.Token)

	err := gocron.Every(15).Minute().Do(weibo.HotSearch)
	if err != nil {
		log.Fatal("go cron error:", err)
	}

	<-gocron.Start()
}
