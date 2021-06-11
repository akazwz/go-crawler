package douban

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"time"
)

func parseDetail(c *colly.Collector, href string) {
	c = c.Clone()
	c.OnHTML("body", func(e *colly.HTMLElement) {
		content := e.DOM.Find("div#content")
		src, _ := content.Find("div#mainpic > a > img").Attr("src")
		rangeNo := content.Find("div.top250 > span.top250-no").Text()
		movieName := content.Find("h1 > span").First().Text()
		year := content.Find("h1 > span.year").Text()

		//director := content.Find("div#info > span").Slice(0, 1).Find("span.attrs > a").Text()
		content.Find("div#info > span").Slice(1, 2).
			Find("span.attrs > a").Each(func(i int, selection *goquery.Selection) {
			screenWriter := selection.Text()
			log.Println("编剧:", screenWriter)
		})

		// 剧情获取有多个span,个数不确定,最少为两个最后一个为 @豆瓣,不需要,remove,取最后一个的全部详情
		content.Find("div#link-report > span").Last().Remove()
		reportSpans := content.Find("div#link-report > span")
		introduction := reportSpans.Last().Text()
		log.Println("\n排名:", rangeNo, "\n电影名:", movieName, "\n海报图片:", src, "\n年份:", year, "\n剧情简介:", introduction)

		ul := content.Find("div#celebrities > ul")
		log.Println("演职员:")
		ul.Find("li").Each(func(i int, selection *goquery.Selection) {
			name := selection.Find("span.name > a").Text()
			role := selection.Find("span.role").Text()
			log.Println("\n姓名:", name, "\n角色:", role)
		})

		rate := content.Find("div.rating_self > strong").Text()
		rateSum := content.Find("div.rating_right > div.rating_sum > a > span").Text()
		log.Println("\n评分:", rate, "\n总评价人数:", rateSum)

	})

	err := c.Visit(href)
	if err != nil {
		log.Fatal("获取详情失败:", err)
	}
}

func Movies() {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("设置频率限制出错:", err)
	}

	// img_movie_details-link 爬取电影详情链接
	c.OnHTML("ol.grid_view", func(e *colly.HTMLElement) {
		// 遍历所有li节点
		e.DOM.Find("li").Each(func(i int, selection *goquery.Selection) {
			href, exists := selection.Find("div.hd > a").Attr("href")
			if exists {
				parseDetail(c, href)
			}
		})
	})

	// img_next_page 爬取下一页链接
	c.OnHTML("div.paginator > span.next", func(e *colly.HTMLElement) {
		href, exists := e.DOM.Find("a").Attr("href")
		log.Println("下一页链接为:", e.Request.AbsoluteURL(href))
		if exists {
			err := e.Request.Visit(e.Request.AbsoluteURL(href))
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	err = c.Visit("https://movie.douban.com/top250")
	if err != nil {
		log.Fatal(err)
	}
}
