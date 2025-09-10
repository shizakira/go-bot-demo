package database

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Password string
}

type Redis struct {
	*redis.Client
}

func NewRedisClient(ctx context.Context, c RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       0,
		Protocol: 2,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &Redis{client}, nil
}
