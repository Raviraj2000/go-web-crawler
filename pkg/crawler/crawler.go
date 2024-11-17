package crawler

import (
	"sync"

	"github.com/Raviraj2000/go-web-crawler/pkg/ratelimiter"
	"github.com/Raviraj2000/go-web-crawler/pkg/redisqueue"
	"github.com/Raviraj2000/go-web-crawler/pkg/storage/models"
	"github.com/Raviraj2000/go-web-crawler/pkg/worker"
)

type Crawler struct {
	Results     chan models.PageData
	WorkerCount int
	RateLimiter *ratelimiter.RateLimiter
	JobCounter  sync.WaitGroup // New field to track active jobs
}

func NewCrawler(workerCount int, rateLimiter *ratelimiter.RateLimiter) *Crawler {
	return &Crawler{
		Results:     make(chan models.PageData, 100),
		WorkerCount: workerCount,
		RateLimiter: rateLimiter,
		JobCounter:  sync.WaitGroup{}, // Initialize the JobCounter WaitGroup
	}
}

func (c *Crawler) Start(redisQueue redisqueue.RedisQueue) {
	for i := 0; i < c.WorkerCount; i++ {
		go worker.Worker(i, c.Results, c.RateLimiter, &c.JobCounter, redisQueue)
	}
}
