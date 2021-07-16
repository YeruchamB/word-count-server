package main

import (
	"github.com/go-redis/redis"
)

var client *redis.Client

// Init redis client
func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

// Increment word count
func Increment(key string, increment int64) error {
	cmd := client.IncrBy(key, increment)
	return cmd.Err()
}

// Get word count
func GetCount(key string) (int64, error) {
	cmd := client.Get(key)
	count, err := cmd.Int64()

	// If err == redis.Nil, the word doesn't exist in the cache and shouldn't return an error
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}
