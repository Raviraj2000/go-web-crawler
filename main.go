package main

import (
	"log"
	"sync"
	"time"

	"github.com/Raviraj2000/go-web-crawler/crawler"
	"github.com/Raviraj2000/go-web-crawler/storage"
)

func main() {
	rateLimiter := crawler.NewRateLimiter(3, 10)

	c := crawler.NewCrawler(10, rateLimiter)
	c.Start()

	seedURLs := []string{
		"https://www.google.com",
		"https://www.golang.org",
		"https://www.nytimes.com",
	}

	for _, url := range seedURLs {
		c.Jobs <- url
	}

	go func() {
		time.Sleep(60 * time.Second)
		close(c.Jobs)
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
