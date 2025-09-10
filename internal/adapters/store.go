package adapters

import (
	"context"
	"encoding/json"
	"github.com/shizakira/daily-tg-bot/pkg/database"
)

type RedisStore[T any] struct {
	client *database.Redis
}

func NewRedisStore[T any](client *database.Redis) *RedisStore[T] {
	return &RedisStore[T]{client: client}
}

func (r *RedisStore[T]) Set(ctx context.Context, key string, value *T) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	r.client.Set(ctx, key, raw, 0)
	return nil
}

func (r *RedisStore[T]) Get(ctx context.Context, key string) (*T, error) {
	raw, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	value := new(T)
	err = json.Unmarshal(raw, value)
	return value, err
}
