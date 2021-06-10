package examples

import (
	"github.com/gocolly/colly"
	"log"
)

func Login() {
	c := colly.NewCollector()

	err := c.Post("https://api.akazwz.com", map[string]string{"username": "akazwz", "password": "123456"})
	if err != nil {
		log.Fatal(err)
	}

	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
	})

	c.Visit("https://api.akazwz.com/")
}
