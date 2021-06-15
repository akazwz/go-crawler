package douban

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"strings"
	"time"
)

func parseDetail(c *colly.Collector, href string) {
	c = c.Clone()
	c.OnHTML("body", func(e *colly.HTMLElement) {
		content := e.DOM.Find("div#content")
		rangeNo := content.Find("div.top250 > span.top250-no").Text()
		movieName := content.Find("h1 > span").First().Text()
		year := content.Find("h1 > span.year").Text()
		moviePost := ""
		imgSrc, exists := content.Find("div#mainpic > a > img").Attr("src")
		if exists {
			moviePost = imgSrc
		}
		// 剧情获取有多个span,个数不确定,最少为两个最后一个为 @豆瓣,不需要,remove,取最后一个的全部详情
		content.Find("div#link-report > span").Last().Remove()
		reportSpans := content.Find("div#link-report > span")
		introduction := reportSpans.Last().Text()
		introduction = strings.ReplaceAll(introduction, "\n", "")
		introduction = strings.ReplaceAll(introduction, "                                                                    　　", "")
		introduction = strings.TrimSpace(introduction)

		// 电影评分以及人数
		rate := content.Find("div.rating_self > strong").Text()
		rateSum := content.Find("div.rating_right > div.rating_sum > a > span").Text()

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

		/**
		排名: No.1
		电影名: 肖申克的救赎 The Shawshank Redemption
		年份: (1994)
		电影海报: https://img2.doubanio.com/view/photo/s_ratio_poster/public/p480747492.jpg
		剧情简介: 一场谋杀案使银行家安迪（蒂姆•罗宾斯 Tim Robbins 饰）蒙冤入狱，谋杀妻子及其情人的指控将囚禁他终生。在肖申克监狱的首次现身就让监狱“大哥”瑞德（摩根•弗里曼 Morgan Freeman 饰）对他另眼相看。瑞德帮助他搞到一把石锤和一幅女明星海报，两人渐成患难 之交。很快，安迪在监狱里大显其才，担当监狱图书管理员，并利用自己的金融知识帮助监狱官避税，引起了典狱长的注意，被招致麾下帮助典狱长洗黑钱。偶然一次，他得知一名新入狱的小偷能够作证帮他洗脱谋杀罪。燃起一丝希望的安迪找到了典狱长，希望他能帮自己翻案。阴险伪善的狱长假装答应安迪，背后却派人杀死小偷，让他唯一能合法出狱的希望泯灭。沮丧的安迪并没有绝望，在一个电闪雷鸣的风雨夜，一场暗藏几十年的越狱计划让他自我救赎，重获自由！老朋友瑞德在他的鼓舞和帮助下，也勇敢地奔向自由。                                                                    　　本片获得1995年奥斯卡10项提名，以及金球奖、土星奖等多项提名。
		导演: 弗兰克·德拉邦特
		编剧: 弗兰克·德拉邦特 / 斯蒂芬·金
		主演: 蒂姆·罗宾斯 / 摩根·弗里曼 / 鲍勃·冈顿 / 威廉姆·赛德勒 / 克兰西·布朗 / 吉尔·贝罗斯 / 马克·罗斯顿 / 詹姆斯·惠特摩 / 杰弗里·德曼 / 拉里·布兰登伯格 / 尼尔·吉恩托利 / 布赖恩·利比 / 大卫·普罗瓦尔 / 约瑟夫·劳格诺 / 祖德·塞克利拉 / 保罗·麦克兰尼 / 芮妮·布莱恩 / 阿方索·弗里曼 / V·J·福斯特 / 弗兰克·梅德拉诺 / 马克·迈尔斯 / 尼尔·萨默斯 / 耐德·巴拉米 / 布赖恩·戴拉特 / 唐·麦克马纳斯
		制片国家/地区: 美国
		语言: 英语
		上映日期: 1994-09-10(多伦多电影节) / 1994-10-14(美国)
		片长: 142分钟
		又名: 月黑高飞(港) / 刺激1995(台) / 地狱诺言 / 铁窗岁月 / 消香克的救赎
		IMDb: tt0111161
		评分: 9.7
		总评价人数: 2370300
		*/

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
