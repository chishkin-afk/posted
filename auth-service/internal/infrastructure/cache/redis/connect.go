package redis

import (
	"context"
	"fmt"

	"github.com/chishkin-afk/posted/auth-service/internal/infrastructure/config"
	"github.com/redis/go-redis/v9"
)

func Connect(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Cache.Redis.Host, cfg.Cache.Redis.Port),
		Username: cfg.Cache.Redis.Auth.Username,
		Password: cfg.Cache.Redis.Auth.Password,
		DB:       cfg.Cache.Redis.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis client: %w", err)
	}

	return client, nil
}

func Close(client *redis.Client) error {
	if err := client.Close(); err != nil {
		return fmt.Errorf("failed to close redis client: %w", err)
	}

	return nil
}
