package crawler

import (
	"sync"

	"github.com/Raviraj2000/go-web-crawler/storage"
)

type Crawler struct {
	Jobs        chan string
	Results     chan storage.PageData
	WorkerCount int
	RateLimiter *RateLimiter
	Visited     map[string]bool
	Mutex       sync.RWMutex
	MaxDepth    int
	JobCounter  sync.WaitGroup // New field to track active jobs
}

func NewCrawler(workerCount int, rateLimiter *RateLimiter) *Crawler {
	return &Crawler{
		Jobs:        make(chan string, 100),
		Results:     make(chan storage.PageData, 100),
		WorkerCount: workerCount,
		RateLimiter: rateLimiter,
		Visited:     make(map[string]bool),
		JobCounter:  sync.WaitGroup{}, // Initialize the JobCounter WaitGroup
	}
}

func (c *Crawler) Start() {
	for i := 0; i < c.WorkerCount; i++ {
		go worker(i, c.Jobs, c.Results, c.RateLimiter, c, &c.JobCounter)
	}
}

func (c *Crawler) MarkVisited(url string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Visited[url] = true
}

func (c *Crawler) IsVisited(url string) bool {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	return c.Visited[url]
}
