package examples

import (
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"os"
)

func generateFormData() map[string][]byte {
	f, _ := os.Open("./README.md")
	defer f.Close()

	fileData, _ := ioutil.ReadAll(f)

	return map[string][]byte{
		"firstname": []byte("Zhao"),
		"lastname":  []byte("WenZhuo"),
		"email":     []byte("akazwz@pm.me"),
		"file":      fileData,
	}
}

func setupServer() {
	var handle http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("received request")
		err := r.ParseMultipartForm(1000000)
		if err != nil {
			fmt.Println("server: Error")
			w.WriteHeader(500)
			w.Write([]byte(
				"<html><body>Internal Server Error</body></html>",
			))
		}
		w.WriteHeader(200)
		fmt.Println("server: OK")
		w.Write([]byte("<html><body>Success</body></html>"))
	}

	go http.ListenAndServe(":8080", handle)
}

func Multipart() {
	setupServer()

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(5),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Posting README.md to", r.URL.String())
	})

	c.PostMultipart("http://localhost:8080/", generateFormData())
	c.Wait()
}
