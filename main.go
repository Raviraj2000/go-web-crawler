package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Raviraj2000/go-web-crawler/crawler"
	"github.com/Raviraj2000/go-web-crawler/redisqueue"
	"github.com/Raviraj2000/go-web-crawler/storage"
	"github.com/go-redis/redis/v8"
)

func main() {
	// Get Redis address from environment
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

	// Get worker count from the environment
	workerCountStr := os.Getenv("WORKER_COUNT")
	workerCount, err := strconv.Atoi(workerCountStr)
	if err != nil || workerCount <= 0 {
		workerCount = 10
		log.Printf("Worker Count not set or invalid. Using default WORKER_COUNT %d\n", workerCount)
	}

	fmt.Printf("Starting web crawler with seed URL: %s and %d workers\n", seedURL, workerCount)

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Initialize RedisQueue with queue and set names for deduplication
	rq := redisqueue.NewRedisQueue(rdb, "url_queue", "visited_urls", seedURL)

	// Add seed URL if unique
	isUnique, err := rq.IsValidURL(seedURL)
	if err != nil {
		log.Fatalf("Error checking or adding seed URL: %v", err)
	}
	if isUnique {
		if err := rq.PushURL(seedURL); err != nil {
			log.Fatalf("Error adding seed URL to queue: %v", err)
		}
		log.Printf("Seed URL added to the queue: %s\n", seedURL)
	}

	// Initialize the crawler with rate limiter and RedisQueue
	rateLimiter := crawler.NewRateLimiter(5, 10)
	c := crawler.NewCrawler(workerCount, rateLimiter)
	c.Start(rq) // Pass RedisQueue instance to the crawler

	// Goroutine to save results and add new URLs to the Redis queue
	go func() {
		for data := range c.Results {
			// Save the crawled data
			storage.Save(data)
			log.Printf("Crawled and saved: %s - %s\n", data.URL, data.Title)
		}
	}()

	// Wait indefinitely (could be improved with graceful shutdown or exit signals)
	select {}
}
