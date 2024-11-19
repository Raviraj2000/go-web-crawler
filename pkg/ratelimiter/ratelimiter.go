package ratelimiter

import (
	"context"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	Limiter *rate.Limiter
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		Limiter: rate.NewLimiter(r, b),
	}
}

func (rl *RateLimiter) Wait() {
	rl.Limiter.Wait(context.Background())
}
