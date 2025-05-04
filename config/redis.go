package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

// InitRedis initializes the Redis connection
func InitRedis() {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	// Test the connection
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Redis connected successfully")
}

// GetRedis returns the Redis client instance
func GetRedis() *redis.Client {
	return RedisClient
}

// SetKey sets a key in Redis with expiration
func SetKey(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(Ctx, key, value, expiration).Err()
}

// GetKey gets a key from Redis
func GetKey(key string) (string, error) {
	return RedisClient.Get(Ctx, key).Result()
}

// DeleteKey deletes a key from Redis
func DeleteKey(key string) error {
	return RedisClient.Del(Ctx, key).Err()
}
