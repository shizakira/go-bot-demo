package redis

import (
	"context"
	"encoding/json"
	"github.com/shizakira/daily-tg-bot/pkg/database"
)

type TaskStateStore[T any] struct {
	client *database.Redis
}

func NewTaskStateStore[T any](client *database.Redis) *TaskStateStore[T] {
	return &TaskStateStore[T]{client: client}
}

func (r *TaskStateStore[T]) Set(ctx context.Context, key string, value *T) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	r.client.Set(ctx, key, raw, 0)
	return nil
}

func (r *TaskStateStore[T]) Get(ctx context.Context, key string) (*T, error) {
	raw, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	value := new(T)
	err = json.Unmarshal(raw, value)
	return value, err
}

func (r *TaskStateStore[T]) Del(ctx context.Context, key string) {
	r.client.Del(ctx, key)
}
