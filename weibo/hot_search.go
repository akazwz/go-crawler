package weibo

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"strconv"
)

func HotSearchDetails(c *colly.Collector, link string, rank int64, content string, hot int64) {
	c = c.Clone()
	// 获取微博热搜详情,得到导语
	c.OnHTML("div#pl_feedlist_index", func(e *colly.HTMLElement) {
		topicLead := e.DOM.Find("div.card-wrap >div.card.card-topic-lead.s-pg16 >p").Text()
		log.Printf("rank: %v, content: %v\n, link: %v\n, hot: %v\n, topic-lead: %v\n",
			rank, content, link, hot, topicLead)
	})

	err := c.Visit(link)
	if err != nil {
		log.Fatal(err)
	}
}

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
			rankStr := selection.Find("td.td-01.ranktop").Text()
			// 判断排名是否为数字,热搜中可能穿插其他内容
			rank, err := strconv.ParseInt(rankStr, 10, 64)
			// 真正热搜内容
			if err == nil {
				// 热搜内容
				content := selection.Find("td.td-02 >a").Text()
				// 热度
				hotStr := selection.Find("td.td-02 >span").Text()
				hot, _ := strconv.ParseInt(hotStr, 10, 64)
				// 热搜链接
				link, exists := selection.Find("td.td-02 >a").Attr("href")
				if exists {
					HotSearchDetails(c, "https://s.weibo.com"+link, rank, content, hot)
				}
			}
		})
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
}
