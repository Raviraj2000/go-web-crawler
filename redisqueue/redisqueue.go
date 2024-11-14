package redisqueue

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
)

type RedisQueue struct {
	Client    *redis.Client
	seedURL   string
	QueueName string
	SetName   string
}

func NewRedisQueue(client *redis.Client, queueName, setName, seedURL string) RedisQueue {
	return RedisQueue{
		Client:    client,
		QueueName: queueName,
		SetName:   setName,
		seedURL:   seedURL,
	}
}

func (rq *RedisQueue) PushURL(url string) error {
	_, err := rq.Client.LPush(context.Background(), rq.QueueName, url).Result()
	return err
}

func (rq *RedisQueue) PopURL() (string, error) {
	url, err := rq.Client.RPop(context.Background(), rq.QueueName).Result()
	if err == redis.Nil {
		return "", nil // No URL in the queue
	}
	return url, err
}

func (rq *RedisQueue) IsValidURL(url string) (bool, error) {
	if !strings.HasPrefix(url, rq.seedURL) {
		return false, nil
	}
	added, err := rq.Client.SAdd(context.Background(), rq.SetName, url).Result()
	if err != nil || added == 0 {
		return false, err
	}
	return true, nil
}
