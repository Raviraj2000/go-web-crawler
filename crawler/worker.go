package crawler

import (
	"log"
	"net/http"

	"github.com/Raviraj2000/go-web-crawler/parser"
	"github.com/Raviraj2000/go-web-crawler/storage"
)

func worker(id int, jobs <-chan string, results chan<- storage.PageData, rateLimiter *RateLimiter, c *Crawler) {
	for url := range jobs {
		c.MarkVisited(url)
		rateLimiter.Wait()
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Worker %d: Error fetching %s", id, err)
			continue
		}
		data, links, err := parser.Parse(resp)
		if err != nil {
			log.Printf("Worker %d: Error parsing %s", id, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		results <- data

		for _, link := range links {
			crawlerEnqueue(c, link)
		}
	}
}

func crawlerEnqueue(c *Crawler, url string) {
	if !c.IsVisited(url) {
		c.Jobs <- url
		return
	}
}
