package infra

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

const RedisKeyNotFound string = "key not found"

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
}

type redisClient struct {
	client *redis.Client
}

// NewRedisClient builds a Redis client
func NewRedisClient() RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	return &redisClient{
		client: client,
	}
}

// Set sets a redis key
func (c *redisClient) Set(ctx context.Context, key string, value string) error {
	if err := c.client.Set(ctx, key, value, 0).Err(); err != nil {
		return errors.New(fmt.Sprintf("error storing value in Redis: %s", err))
	}

	return nil
}

// Get gets a redis key
func (c *redisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New(RedisKeyNotFound)
	}

	if err != nil {
		return "", errors.New(fmt.Sprintf("error getting value from Redis: %s", err))
	}

	return value, nil
}
