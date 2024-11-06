package crawler

import (
	"log"
	"net/http"
	"sync"

	"github.com/Raviraj2000/go-web-crawler/parser"
	"github.com/Raviraj2000/go-web-crawler/storage"
)

func worker(id int, jobs <-chan string, results chan<- storage.PageData, rateLimiter *RateLimiter, c *Crawler, jobCounter *sync.WaitGroup) {
	for url := range jobs {
		c.MarkVisited(url)
		rateLimiter.Wait()
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Worker %d: Error fetching %s", id, err)
			jobCounter.Done() // Decrement if job fails
			continue
		}

		data, links, err := parser.Parse(resp)
		if err != nil {
			log.Printf("Worker %d: Error parsing %s", id, err)
			resp.Body.Close()
			jobCounter.Done() // Decrement if parsing fails
			continue
		}
		resp.Body.Close()
		results <- data

		for _, link := range links {
			crawlerEnqueue(c, link, jobCounter) // Enqueue each link
		}
		jobCounter.Done() // Decrement after processing all links
	}
}

func crawlerEnqueue(c *Crawler, url string, jobCounter *sync.WaitGroup) {
	if !c.IsVisited(url) {
		select {
		case c.Jobs <- url:
			jobCounter.Add(1) // Increment job counter for each new URL added
			log.Println("Enqueued URL:", url)
		default:
			log.Printf("Failed to enqueue URL: %s - channel may be closed", url)
		}
	}
}
