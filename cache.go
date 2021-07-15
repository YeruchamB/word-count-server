package main

import (
	"github.com/go-redis/redis"
	"strings"
)

var client *redis.Client

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func Increment(key string, increment int64) error {
	cmd := client.IncrBy(strings.ToLower(key), increment)
	return cmd.Err()
}

func GetCount(key string) (int64, error) {
	cmd := client.Get(strings.ToLower(key))
	count, err := cmd.Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}
