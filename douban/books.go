package douban

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"strings"
	"time"
)

func parseBookDetail(c *colly.Collector, href string) {
	c = c.Clone()

	c.OnHTML("body", func(e *colly.HTMLElement) {
		wrapper := e.DOM.Find("div#wrapper")
		rangeNo := wrapper.Find("div.rank-label > span.rank-label-no").Text()
		rangeNo = strings.TrimSpace(rangeNo) // 排名
		bookName := wrapper.Find("h1").Text()
		bookName = strings.TrimSpace(bookName) // 书名

		bookImg := ""
		content := wrapper.Find("div#content")
		src, exists := content.Find("div#mainpic > a").Attr("href")
		if exists {
			bookImg = src
		}

		_ = content.Find("div#info").Find("span.pl").First().NextAll().Text()
		content.Find("div#info").Find("span").First().Remove()
		infoText := content.Find("div#info").Text()
		infoText = strings.TrimSpace(infoText)
		infoText = strings.ReplaceAll(infoText, "\n", "")

		// 用空格分割字符串
		split := strings.Split(infoText, "              ")

		// 定义一个新切片
		var info = make([]string, 0)

		for _, s := range split {
			s = strings.TrimSpace(s)
			// 去除空格后,长度大于一的放入新切片,得到图书信息
			if len(s) > 1 {
				info = append(info, s)
			}
		}
		for _, s := range info {
			log.Println("this is s", s)
		}

		log.Println(
			"\n排名:", rangeNo,
			"\n书名:", bookName,
			"\n图书封面:", bookImg,
		)
	})

	err := c.Visit(href)
	if err != nil {
		log.Fatal(err, "详情爬取失败")
	}
}

func Books() {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36 Edg/91.0.864.48"
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("设置频率限制出错")
	}

	c.OnHTML("div#content", func(e *colly.HTMLElement) {
		e.DOM.Find("div.article > div.indent > table").Each(func(i int, selection *goquery.Selection) {
			td := selection.Find("tbody > tr > td").Last()
			href, exists := td.Find("div.pl2 > a").Attr("href")
			if exists {
				parseBookDetail(c, href)
			}

			rate := td.Find("div.star > span.rating_nums").Text()
			rateSum := td.Find("div.star > span.pl").Text()
			rateSum = strings.ReplaceAll(rateSum, "(", "")
			rateSum = strings.ReplaceAll(rateSum, ")", "")
			rateSum = strings.ReplaceAll(rateSum, "人评价", "")
			rateSum = strings.TrimSpace(rateSum)
			quote := td.Find("p.quote >span").Text()
			log.Println(
				"\n评分:", rate,
				"\n评分人数:", rateSum,
				"\n引用:", quote,
				"\n--------------------------------------------------",
			)

		})
	})

	c.OnHTML("div.paginator > span.next", func(e *colly.HTMLElement) {
		href, exists := e.DOM.Find("a").Attr("href")
		if exists {
			err := e.Request.Visit(href)
			if err != nil {
				log.Fatal("爬取下一页失败")
			}
		}
	})

	err = c.Visit("https://book.douban.com/top250")
	if err != nil {
		log.Fatal(err)
	}

}
