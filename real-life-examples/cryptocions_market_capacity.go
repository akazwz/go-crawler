package real_life_examples

import (
	"encoding/csv"
	"github.com/gocolly/colly"
	"log"
	"os"
)

func CryptocoinsMarketCapacity() {
	fName := "cryptocionsmarketcap.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{
		"Name", "Symbol", "Price (USD)", "Volume (USD)", "Market capacity (USD)",
		"Change (1h)", "Change (24h)", "Change (7d)",
	})
	if err != nil {
		log.Fatal(err)
	}

	c := colly.NewCollector()

	c.OnHTML("#currencies-all tbody tr", func(e *colly.HTMLElement) {
		err := writer.Write([]string{
			e.ChildText(".currency-name-container"),
			e.ChildText(".col-symbol"),
			e.ChildAttr("a.price", "data-usd"),
			e.ChildAttr("a.volume", "data-usd"),
			e.ChildAttr(".market-cap", "data-usd"),
			e.ChildText(".percent-1h"),
			e.ChildText(".percent-24h"),
			e.ChildText(".percent-7d"),
		})
		if err != nil {
			log.Fatal(err)
		}
	})

	err = c.Visit("https://coinmarketcap.com/all/views/all/")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
