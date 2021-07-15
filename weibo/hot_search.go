package weibo

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"strconv"
)

func HotSearch() {
	c := colly.NewCollector()
	url := "https://s.weibo.com/top/summary/"

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting: ", r.URL)
	})

	c.OnHTML("div#pl_top_realtimehot >table >tbody", func(e *colly.HTMLElement) {
		e.DOM.Find("tr").Each(func(i int, selection *goquery.Selection) {
			if i == 0 {
				// 不是真正热搜内容
			}
			// 热搜排名
			rank := selection.Find("td.td-01.ranktop").Text()
			// 判断排名是否为数子,热搜中可能穿插其他内容
			_, err := strconv.ParseFloat(rank, 64)
			// 真正热搜内容
			if err == nil {
				// 热搜内容
				content := selection.Find("td.td-02 >a").Text()
				// 热度
				hot := selection.Find("td.td-02 >span").Text()
				// 热搜链接
				link, _ := selection.Find("td.td-02 >a").Attr("href")
				log.Printf("rank: %v, content: %v, link: %v, hot: %v",
					rank, content, link, hot)
			}
		})
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
}
