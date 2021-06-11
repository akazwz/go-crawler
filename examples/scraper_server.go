package examples

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"log"
	"net/http"
)

type pageInfo struct {
	StatusCode int
	Links      map[string]int
}

func handler(w http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}

	log.Println("visiting", URL)

	c := colly.NewCollector()

	p := &pageInfo{Links: make(map[string]int)}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link != "" {
			p.Links[link]++
		}
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		p.StatusCode = r.StatusCode
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	c.Visit(URL)

	b, err := json.Marshal(p)
	if err != nil {
		log.Println("failed to serialize response:", err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func ScraperServer() {
	addr := ":7171"

	http.HandleFunc("/", handler)

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
