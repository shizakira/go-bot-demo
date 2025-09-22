package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type TelegramSession struct {
	client *Redis
	ttl    time.Duration
}

func NewTelegramSession(client *Redis) *TelegramSession {
	return &TelegramSession{
		client: client,
		ttl:    24 * time.Hour,
	}
}

func (t *TelegramSession) InitSession(ctx context.Context, chatId string) error {
	exist, err := t.client.Exists(ctx, chatId).Result()
	if err != nil {
		return fmt.Errorf("failed to check session existence: %w", err)
	}
	if exist == 0 {
		err := t.client.HSet(ctx, chatId, "initialized", "true").Err()
		if err != nil {
			return fmt.Errorf("failed to initialize session: %w", err)
		}
	}
	err = t.client.Expire(ctx, chatId, t.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set session TTL: %w", err)
	}
	return nil
}

func (t *TelegramSession) Set(ctx context.Context, chatId string, key string, value any) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	err = t.client.HSet(ctx, chatId, key, raw).Err()
	if err != nil {
		return fmt.Errorf("failed to set session value: %w", err)
	}
	t.client.Expire(ctx, chatId, t.ttl)
	return nil
}

func (t *TelegramSession) Get(ctx context.Context, chatId string, key string) ([]byte, error) {
	raw, err := t.client.HGet(ctx, chatId, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get session value: %w", err)
	}
	t.client.Expire(ctx, chatId, t.ttl)
	return raw, nil
}

func (t *TelegramSession) Del(ctx context.Context, chatId string, key string) error {
	err := t.client.HDel(ctx, chatId, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session value: %w", err)
	}
	return nil
}

func (t *TelegramSession) Clear(ctx context.Context, chatId string) error {
	err := t.client.Del(ctx, chatId).Err()
	if err != nil {
		return fmt.Errorf("failed to clear session: %w", err)
	}
	return nil
}
