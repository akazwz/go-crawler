package douban

import (
	"github.com/PuerkitoBio/goquery"
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
		info := content.Find("div#info")
		infoText := info.Text()
		director := processInfo(infoText, "导演")
		scriptWriter := processInfo(infoText, "编剧")
		mainActors := processInfo(infoText, "主演")
		region := processInfo(infoText, "制片国家/地区")
		languages := processInfo(infoText, "语言")
		date := processInfo(infoText, "上映日期")
		long := processInfo(infoText, "片长")
		anotherName := processInfo(infoText, "又名")
		IMDb := processInfo(infoText, "IMDb")
		log.Println(
			"\n排名:", rangeNo, "\n电影名:", movieName, "\n年份:", year, "\n电影海报:", moviePost, "\n剧情简介:", introduction,
			"\n导演:", director, "\n编剧:", scriptWriter, "\n主演:", mainActors,
			"\n制片国家/地区:", region, "\n语言:", languages, "\n上映日期:", date, "\n片长:", long,
			"\n又名:", anotherName, "\nIMDb:", IMDb,
			"\n评分:", rate, "\n总评价人数:", rateSum,
			"\n---------------------------------------------------------",
		)

		// 获取演职员
		/*ul := content.Find("div#celebrities > ul")
		log.Println("演职员:")
		ul.Find("li").Each(func(i int, selection *goquery.Selection) {
			name := selection.Find("span.name > a").Text()
			role := selection.Find("span.role").Text()
			log.Println("\n姓名:", name, "\n角色:", role)
		})*/
	})

	err := c.Visit(href)
	if err != nil {
		log.Fatal("获取详情失败:", err)
	}
}

// 处理电影详情,传入电影详情text,和需要获取的名称,比如导演,返回相对应的值
func processInfo(infoText, key string) string {
	split := strings.Split(infoText, ":")
	var splitStr string
	switch key {
	case "导演":
		key = "编剧"
		splitStr = split[1]
	case "编剧":
		key = "主演"
		splitStr = split[2]
	case "主演":
		key = "类型"
		splitStr = split[3]
	case "类型":
		key = "制片国家/地区"
		splitStr = split[4]
	case "制片国家/地区":
		key = "语言"
		splitStr = split[5]
	case "语言":
		key = "上映日期"
		splitStr = split[6]
	case "上映日期":
		key = "片长"
		splitStr = split[7]
	case "片长":
		key = "又名"
		splitStr = split[8]
	case "又名":
		key = "IMDb"
		splitStr = split[9]
	case "IMDb":
		key = ""
		splitStr = split[10]
	default:
		log.Fatal("key有误!")
	}

	str := strings.ReplaceAll(splitStr, key, "")
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.TrimSpace(str)
	return str
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
