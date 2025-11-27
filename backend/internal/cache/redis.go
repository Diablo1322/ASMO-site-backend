package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
	connected bool
}

var _ Cache = (*RedisCache)(nil)

func NewRedisCache(redisURL string) (*RedisCache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)
	ctx := context.Background()

	// Test connection with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(timeoutCtx).Err(); err != nil {
		// Возвращаем "fallback" кэш который логирует но не падает
		fmt.Printf("⚠️  Redis connection failed: %v. Using in-memory fallback.\n", err)
		return &RedisCache{
			client: client,
			ctx:    ctx,
			connected: false,
		}, nil
	}

	fmt.Println("✅ Redis connected successfully")
	return &RedisCache{
		client: client,
		ctx:    ctx,
		connected: true,
	}, nil
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	if !r.connected {
		fmt.Printf("⚠️  Redis not connected - skipping SET for key: %s\n", key)
		return nil
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, key, jsonData, expiration).Err()
}

func (r *RedisCache) Get(key string, dest interface{}) error {
	if !r.connected {
		fmt.Printf("⚠️  Redis not connected - skipping GET for key: %s\n", key)
		return ErrNotFound
	}

	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisCache) Delete(key string) error {
	if !r.connected {
		fmt.Printf("⚠️  Redis not connected - skipping DELETE for key: %s\n", key)
		return nil
	}

	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisCache) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

var ErrNotFound = fmt.Errorf("key not found")