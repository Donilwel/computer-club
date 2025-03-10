package repository

import (
	"computer-club/internal/config"
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr(),
		DB:   0,
	})

	// Проверяем соединение
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	return client
}
