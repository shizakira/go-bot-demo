package telegram

import "context"

type Session interface {
	InitSession(ctx context.Context, chatId string) error
	Set(ctx context.Context, chatId string, key string, value any) error
	Get(ctx context.Context, chatId string, key string) ([]byte, error)
	Del(ctx context.Context, chatId string, key string) error
	Clear(ctx context.Context, chatId string) error
}
