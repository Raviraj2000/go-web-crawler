package worker

import (
	"log"
	"net/http"
	"sync"

	"github.com/Raviraj2000/go-web-crawler/pkg/parser"
	"github.com/Raviraj2000/go-web-crawler/pkg/ratelimiter"
	"github.com/Raviraj2000/go-web-crawler/pkg/redisqueue"
	"github.com/Raviraj2000/go-web-crawler/pkg/storage/models"
)

func Worker(id int, results chan<- models.PageData, rateLimiter *ratelimiter.RateLimiter, jobCounter *sync.WaitGroup, redisQueue redisqueue.RedisQueue) {
	for {
		// Fetch URL from Redis queue
		url, err := redisQueue.PopURL()
		if err != nil {
			log.Printf("Worker %d: Error fetching URL from Redis: %v", id, err)
			continue
		}
		if url == "" {
			// If no URL is found, wait and retry
			log.Printf("Worker %d: No URLs in queue, retrying...", id)
			continue
		}
		rateLimiter.Wait()
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Worker %d: Error fetching %s: %v", id, url, err)
			continue
		}

		data, links, err := parser.Parse(resp)
		resp.Body.Close()
		if err != nil {
			log.Printf("Worker %d: Error parsing %s: %v", id, url, err)
			continue
		}

		results <- data

		for _, link := range links {
			crawlerEnqueue(link, jobCounter, redisQueue) // Enqueue each link in Redis
		}
		jobCounter.Done() // Decrement after processing all links
	}
}

func crawlerEnqueue(url string, jobCounter *sync.WaitGroup, redisQueue redisqueue.RedisQueue) {
	// Check if URL is unique in Redis before adding to queue
	isValid, err := redisQueue.IsValidURL(url)
	if err != nil {
		log.Printf("Error checking URL uniqueness: %v", err)
		return
	}
	if isValid {
		err := redisQueue.PushURL(url)
		if err != nil {
			log.Printf("Failed to enqueue URL in Redis: %s - %v", url, err)
		} else {
			jobCounter.Add(1) // Increment job counter for each new URL added
			log.Println("Enqueued URL:", url)
		}
	}
}
