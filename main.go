package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Raviraj2000/go-web-crawler/crawler"
	"github.com/Raviraj2000/go-web-crawler/storage"
)

func main() {
	// Get environment variable for seed URL
	seedURL := os.Getenv("SEED_URL")
	if seedURL == "" {
		fmt.Println("SEED_URL not provided. Please enter a seed URL:")
		reader := bufio.NewReader(os.Stdin)
		seedURL, _ = reader.ReadString('\n')
		seedURL = strings.TrimSpace(seedURL) // Remove any newline or extra spaces
	}

	// Stop the program if no URL is provided
	if seedURL == "" {
		fmt.Println("No seed URL provided. Exiting program.")
		return
	}

	// Get environment variable for worker count
	workerCountStr := os.Getenv("WORKER_COUNT")
	workerCount, err := strconv.Atoi(workerCountStr)
	if err != nil || workerCount <= 0 {
		workerCount = 1000 // Default worker count if not provided or invalid
	}

	// Initialize rate limiter and crawler
	rateLimiter := crawler.NewRateLimiter(5, 10)
	c := crawler.NewCrawler(workerCount, rateLimiter)
	c.Start()

	// Add the seed URL to the job queue
	c.JobCounter.Add(1)
	c.Jobs <- seedURL

	// Close the Jobs channel when all URLs are processed
	go func() {
		c.JobCounter.Wait() // Wait until all jobs are done
		close(c.Jobs)       // Close Jobs only when no more URLs are left to enqueue
	}()

	// Process results and log crawled data
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
