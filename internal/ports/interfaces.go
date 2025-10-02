package ports

import (
	"context"
	"github.com/shizakira/daily-tg-bot/internal/domain"
)

type TaskRepository interface {
	Add(ctx context.Context, t domain.Task) error
	GetOpenByUserID(ctx context.Context, userUd int64) ([]*domain.Task, error)
	CloseTask(ctx context.Context, id int64, isDone bool) error
	GetExpiredTasks(ctx context.Context) ([]*domain.Task, error)
	GetSoonExpiredTasks(ctx context.Context) ([]*domain.Task, error)
}

type TelegramUserRepository interface {
	Create(ctx context.Context, user *domain.TelegramUser) error
	FindByChatID(ctx context.Context, chatId int64) (*domain.TelegramUser, error)
	FindByUserIDs(ctx context.Context, userIDs []int64) ([]*domain.TelegramUser, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (int64, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	Exists(ctx context.Context, id int64) (bool, error)
}
