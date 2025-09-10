package telegram

import (
	"context"
)

type TaskStateStore interface {
	Set(ctx context.Context, key string, value *TaskState) error
	Get(ctx context.Context, key string) (*TaskState, error)
}
