package weibo

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/akazwz/go-crawler/utils/influx"
	"github.com/gocolly/colly"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HotSearchDetails(c *colly.Collector, link string, rank string, content string, hot string, t time.Time) {
	log.Println("开始爬取详情")
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
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	url := "https://s.weibo.com/top/summary/"

	tags := map[string]string{}
	fields := map[string]interface{}{}

	tags["rank"] = "00"
	fields["pdf_file"] = "pdf"
	fields["image_file"] = "image"

	// 一次热搜应该为同一时间
	t := time.Now()
	err := influx.Write("hot_search", tags, fields, t)
	if err != nil {
		log.Println("influx error:", err)
	}

	c.OnRequest(func(r *colly.Request) {
	})

	log.Println("开始")

	c.OnHTML("body", func(e *colly.HTMLElement) {
		log.Println("开始爬取热搜")
		e.DOM.Find("tr").Each(func(i int, selection *goquery.Selection) {
			log.Println("找到热搜内容")
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
				/*keywords := make([]string, 2)
				keywords[0] = "iPhone"
				keywords[1] = "姐姐"
				for i := 0; i < len(keywords); i++ {
					if strings.Contains(content, keywords[i]) {
						log.Println(keywords[i], "上热搜了")
						//notify.SendMessage()
					}
				}*/
				// 热度
				hot := selection.Find("td.td-02 >span").Text()
				// 热搜链接
				link, exists := selection.Find("td.td-02 >a").Attr("href")
				log.Println("link:", link)
				if exists {
					// 以 / 开头的为有效链接
					if strings.HasPrefix(link, "/") {
						log.Println("有效链接:", "https://s.weibo.com"+link)
						HotSearchDetails(c, "https://s.weibo.com"+link, rank, content, hot, t)
					}
				}
			}
		})
	})

	err = c.Visit(url)
	if err != nil {
		log.Println(err)
	}
}
