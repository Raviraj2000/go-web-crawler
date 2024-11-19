package crawler_test

import (
	"testing"

	"github.com/Raviraj2000/go-web-crawler/pkg/crawler"
	"github.com/Raviraj2000/go-web-crawler/pkg/ratelimiter"
)

func TestNewCrawler(t *testing.T) {
	workerCount := 5
	rateLimiter := ratelimiter.NewRateLimiter(10, 5)

	c := crawler.NewCrawler(workerCount, rateLimiter)

	if c.WorkerCount != workerCount {
		t.Errorf("Expected WorkerCount %d, got %d", workerCount, c.WorkerCount)
	}

	if c.RateLimiter != rateLimiter {
		t.Error("RateLimiter instance does not match")
	}

	if c.Results == nil {
		t.Error("Results channel is not initialized")
	}
}
