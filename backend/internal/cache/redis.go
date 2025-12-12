package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
        log.Printf("⚠️ Failed to parse Redis URL: %v", err)
        return &RedisCache{connected: false}, nil  // Fallback без паники
    }

    client := redis.NewClient(opts)

    // Таймаут подключения
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.Ping(ctx).Err(); err != nil {
        log.Printf("⚠️ Redis connection failed: %v. Using fallback mode.", err)
        return &RedisCache{
            client: client,
            ctx:    context.Background(),
            connected: false,
        }, nil
    }

    return &RedisCache{
        client: client,
        ctx:    context.Background(),
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