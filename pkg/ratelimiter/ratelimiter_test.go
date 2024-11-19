package ratelimiter_test

import (
	"testing"
	"time"

	"github.com/Raviraj2000/go-web-crawler/pkg/ratelimiter"
	"golang.org/x/time/rate"
)

func TestNewRateLimiter(t *testing.T) {
	r := rate.Every(100 * time.Millisecond)
	b := 5
	rl := ratelimiter.NewRateLimiter(r, b)

	if rl == nil {
		t.Fatalf("Expected a RateLimiter instance, got nil")
	}

	if rl.Limiter == nil {
		t.Fatalf("Expected a rate.Limiter instance inside RateLimiter, got nil")
	}

	if rl.Limiter.Limit() != r {
		t.Errorf("Expected rate limit to be %v, got %v", r, rl.Limiter.Limit())
	}

	if rl.Limiter.Burst() != b {
		t.Errorf("Expected burst size to be %d, got %d", b, rl.Limiter.Burst())
	}
}

func TestRateLimiterWait(t *testing.T) {
	r := rate.Every(50 * time.Millisecond)
	b := 1
	rl := ratelimiter.NewRateLimiter(r, b)

	start := time.Now()

	// First call should pass immediately as the bucket has a token
	rl.Wait()

	// Second call should wait for the rate duration
	rl.Wait()

	elapsed := time.Since(start)
	expectedWait := 50 * time.Millisecond

	// Allow some leeway for timing variations
	if elapsed < expectedWait {
		t.Errorf("Expected to wait at least %v, but waited %v", expectedWait, elapsed)
	}
}
