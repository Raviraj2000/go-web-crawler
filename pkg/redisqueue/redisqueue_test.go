package redisqueue_test

import (
	"context"
	"testing"

	"github.com/Raviraj2000/go-web-crawler/pkg/redisqueue"
	"github.com/go-redis/redis/v8"
)

func setupTestRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0, // Use default DB
	})
	// Flush the database to ensure a clean state for testing
	client.FlushDB(context.Background())
	return client
}

func TestPushAndPopURL(t *testing.T) {
	client := setupTestRedis()
	defer client.Close()

	rq := redisqueue.NewRedisQueue(client, "test_queue", "test_set", "https://example.com")

	// Test pushing a URL
	url := "https://example.com/page1"
	err := rq.PushURL(url)
	if err != nil {
		t.Fatalf("Failed to push URL: %v", err)
	}

	// Test popping a URL
	poppedURL, err := rq.PopURL()
	if err != nil {
		t.Fatalf("Failed to pop URL: %v", err)
	}
	if poppedURL != url {
		t.Errorf("Expected %q, got %q", url, poppedURL)
	}

	// Test popping from an empty queue
	poppedURL, err = rq.PopURL()
	if err != nil {
		t.Fatalf("Error popping from an empty queue: %v", err)
	}
	if poppedURL != "" {
		t.Errorf("Expected empty string, got %q", poppedURL)
	}
}

func TestIsValidURL(t *testing.T) {
	client := setupTestRedis()
	defer client.Close()

	rq := redisqueue.NewRedisQueue(client, "test_queue", "test_set", "https://example.com")

	// Test adding a valid URL
	url := "https://example.com/page1"
	isValid, err := rq.IsValidURL(url)
	if err != nil {
		t.Fatalf("Error checking URL validity: %v", err)
	}
	if !isValid {
		t.Errorf("Expected URL %q to be valid", url)
	}

	// Test adding the same URL again (should be invalid)
	isValid, err = rq.IsValidURL(url)
	if err != nil {
		t.Fatalf("Error checking URL validity: %v", err)
	}
	if isValid {
		t.Errorf("Expected URL %q to be invalid on duplicate addition", url)
	}

	// Test adding a URL outside the seed domain
	invalidURL := "https://other.com/page"
	isValid, err = rq.IsValidURL(invalidURL)
	if err != nil {
		t.Fatalf("Error checking URL validity: %v", err)
	}
	if isValid {
		t.Errorf("Expected URL %q to be invalid", invalidURL)
	}
}
