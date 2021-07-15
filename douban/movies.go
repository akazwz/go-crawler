package douban

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/akazwz/go-crawler/utils"
	"github.com/gocolly/colly"
	"log"
	"strings"
	"time"
)

func parseMovieDetail(c *colly.Collector, href string) {
	c = c.Clone()
	c.OnHTML("body", func(e *colly.HTMLElement) {
		content := e.DOM.Find("div#content")
		rangeNo := content.Find("div.top250 > span.top250-no").Text() // 排名编号
		movieName := content.Find("h1 > span").First().Text()         // 电影名
		year := content.Find("h1 > span.year").Text()                 //年份
		year = strings.ReplaceAll(year, "(", "")
		year = strings.ReplaceAll(year, ")", "")
		moviePost := ""
		imgSrc, exists := content.Find("div#mainpic > a > img").Attr("src")
		if exists {
			moviePost = imgSrc // 电影海报
		}
		// 剧情获取有多个span,个数不确定,最少为两个最后一个为 @豆瓣,不需要,remove,取最后一个的全部详情
		content.Find("div#link-report > span").Last().Remove()
		reportSpans := content.Find("div#link-report > span")
		introduction := reportSpans.Last().Text() // 电影简介
		introduction = strings.ReplaceAll(introduction, "\n", "")
		introduction = strings.ReplaceAll(introduction, "                                                                    　　", "")
		introduction = strings.TrimSpace(introduction)

		// 电影评分以及人数
		rate := content.Find("div.rating_self > strong").Text()                        // 电影评分
		rateSum := content.Find("div.rating_right > div.rating_sum > a > span").Text() // 电影评价人数

		// 电影详情
		/*info := content.Find("div#info")
		infoText := info.Text()*/
		record := []string{rangeNo, movieName, moviePost, rate, rateSum}
		utils.WriteCSV("movie-top-250.csv", record)
	})

	err := c.Visit(href)
	if err != nil {
		log.Fatal("获取详情失败:", err)
	}
}

func Movies() {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.67"
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 10 * time.Second,
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
				parseMovieDetail(c, href)
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
