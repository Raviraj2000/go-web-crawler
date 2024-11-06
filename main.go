package main

import (
	"log"
	"sync"

	"github.com/Raviraj2000/go-web-crawler/crawler"
	"github.com/Raviraj2000/go-web-crawler/storage"
)

func main() {
	rateLimiter := crawler.NewRateLimiter(3, 10)
	c := crawler.NewCrawler(1000, rateLimiter)
	c.Start()

	seedURLs := []string{
		"https://www.google.com",
		"https://www.golang.org",
		"https://www.nytimes.com",
	}

	for _, url := range seedURLs {
		c.JobCounter.Add(1) // Increment for each seed URL
		c.Jobs <- url
	}

	// Close the Jobs channel when all URLs are processed
	go func() {
		c.JobCounter.Wait() // Wait until all jobs are done
		close(c.Jobs)       // Close Jobs only when no more URLs are left to enqueue
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range c.Results {
			storage.Save(data)
			log.Printf("Crawled: %s - %s\n", data.URL, data.Title)
		}
	}()

	wg.Wait()
}
