package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, Colly!")

	/*global.VP = initialize.InitViper()
	if global.VP == nil {
		fmt.Println("配置文件初始化失败")
	}

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatal("时区加载失败")
	}

	c := cron.New(cron.WithLocation(location))
	_, err = c.AddFunc("0,15,30,45 * * * * ", weibo.HotSearch)
	if err != nil {
		log.Fatal("定时任务添加失败", err)
	}
	c.Run()
	c.Start()*/

	c := fmt.Sprintf("%02d", 1)
	fmt.Println(c)
}
