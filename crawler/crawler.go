package crawler

import (
	"sync"

	"github.com/Raviraj2000/go-web-crawler/redisqueue"
	"github.com/Raviraj2000/go-web-crawler/storage"
)

type Crawler struct {
	Results     chan storage.PageData
	WorkerCount int
	RateLimiter *RateLimiter
	JobCounter  sync.WaitGroup // New field to track active jobs
}

func NewCrawler(workerCount int, rateLimiter *RateLimiter) *Crawler {
	return &Crawler{
		Results:     make(chan storage.PageData, 100),
		WorkerCount: workerCount,
		RateLimiter: rateLimiter,
		JobCounter:  sync.WaitGroup{}, // Initialize the JobCounter WaitGroup
	}
}

func (c *Crawler) Start(redisQueue redisqueue.RedisQueue) {
	for i := 0; i < c.WorkerCount; i++ {
		go worker(i, c.Results, c.RateLimiter, &c.JobCounter, redisQueue)
	}
}
