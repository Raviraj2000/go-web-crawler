package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Raviraj2000/go-web-crawler/crawler"
	"github.com/Raviraj2000/go-web-crawler/storage"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	// Get Redis address from the environment
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR environment variable not set")
		return
	}

	// Get the seed URL from the environment
	seedURL := os.Getenv("SEED_URL")
	if seedURL == "" {
		log.Fatal("SEED_URL environment variable not set")
		return
	}

	// Get the worker count from the environment
	workerCountStr := os.Getenv("WORKER_COUNT")
	workerCount, err := strconv.Atoi(workerCountStr)
	if err != nil || workerCount <= 0 {
		workerCount = 10 // Default worker count if not provided or invalid
		log.Printf("Worker Count not set or invalid. Using default WORKER_COUNT %d\n", workerCount)
	}

	// Set up Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Push the initial seed URL if the queue is empty
	queueLen, err := rdb.LLen(ctx, "url_queue").Result()
	if err != nil {
		log.Fatalf("Failed to check Redis queue: %v", err)
	}
	if queueLen == 0 {
		if _, err := rdb.LPush(ctx, "url_queue", seedURL).Result(); err != nil {
			log.Fatalf("Failed to push seed URL: %v", err)
		}
		log.Printf("Seed URL added to the queue: %s\n", seedURL)
	}

	// Initialize rate limiter and crawler
	rateLimiter := crawler.NewRateLimiter(5, 10)
	c := crawler.NewCrawler(workerCount, rateLimiter)
	c.Start()

	// Wait group to manage concurrent goroutines
	var wg sync.WaitGroup
	wg.Add(1)

	// Goroutine to fetch URLs from Redis queue and process them
	go func() {
		defer wg.Done()
		for {
			// Fetch URL from the Redis queue
			url, err := rdb.RPop(ctx, "url_queue").Result()
			if err == redis.Nil {
				// Queue is empty, wait a bit and try again
				time.Sleep(1 * time.Second)
				continue
			} else if err != nil {
				log.Printf("Error fetching URL from queue: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// Process the URL with the crawler
			c.JobCounter.Add(1)
			c.Jobs <- url
		}
	}()

	// Goroutine to save results and add new URLs to the queue
	go func() {
		for data := range c.Results {
			storage.Save(data)
			log.Printf("Crawled: %s - %s\n", data.URL, data.Title)
			for _, link := range data.URL {
				// Push new URLs to the Redis queue
				if _, err := rdb.LPush(ctx, "url_queue", link).Result(); err != nil {
					log.Printf("Error pushing URL to queue: %v", err)
				}
			}
		}
	}()

	// Wait for the main goroutine to finish processing
	wg.Wait()
}
