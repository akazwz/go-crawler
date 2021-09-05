package weibo

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/akazwz/go-crawler/utils/influx"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"time"
)

func HotSearchDetails(c *colly.Collector, link string, rank string, content string, hot string, t time.Time) {
	c = c.Clone()
	// 获取微博热搜详情,得到导语
	c.OnHTML("div#pl_feedlist_index", func(e *colly.HTMLElement) {
		topicLead := e.DOM.Find("div.card-wrap >div.card.card-topic-lead.s-pg16 >p").Text()

		tags := map[string]string{}
		fields := map[string]interface{}{}
		rankInt, err := strconv.Atoi(rank)
		if err != nil {
			log.Fatal("string to int error")
		}
		tags["rank"] = fmt.Sprintf("%02d", rankInt)
		fields["content"] = content
		fields["hot"] = hot
		fields["topic_lead"] = topicLead
		fields["link"] = link

		err = influx.Write("hot_search", tags, fields, t)
		if err != nil {
			log.Fatal("influx error:", err)
		}
	})

	err := c.Visit(link)
	if err != nil {
		log.Fatal(err)
	}
}

func HotSearch() {
	c := colly.NewCollector()
	url := "https://s.weibo.com/top/summary/"
	// 一次热搜应该为同一时间
	t := time.Now()
	c.OnRequest(func(r *colly.Request) {
	})

	c.OnHTML("div#pl_top_realtimehot >table >tbody", func(e *colly.HTMLElement) {
		e.DOM.Find("tr").Each(func(i int, selection *goquery.Selection) {
			if i == 0 {
				// 不是真正热搜内容
			}
			// 热搜排名
			rank := selection.Find("td.td-01.ranktop").Text()
			// 判断排名是否为数字,热搜中可能穿插其他内容
			_, err := strconv.ParseInt(rank, 10, 64)
			// 真正热搜内容
			if err == nil {
				// 热搜内容
				content := selection.Find("td.td-02 >a").Text()
				// 热度
				hot := selection.Find("td.td-02 >span").Text()
				// 热搜链接
				link, exists := selection.Find("td.td-02 >a").Attr("href")
				if exists {
					HotSearchDetails(c, "https://s.weibo.com"+link, rank, content, hot, t)
				}
			}
		})
	})

	err := c.Visit(url)
	if err != nil {
		log.Println(err)
	}
}
