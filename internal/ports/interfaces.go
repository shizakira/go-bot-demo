package ports

import (
	"context"

	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/internal/dto"
)

type TaskRepository interface {
	Add(ctx context.Context, t domain.Task) error
	GetOpenByUserID(ctx context.Context, userID int64) ([]*domain.Task, error)
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

type TaskService interface {
	CreateTask(ctx context.Context, input dto.CreateTaskInput) error
	GetOpenTasksByUserID(ctx context.Context, input dto.GetAllTasksByUserIdInput) (dto.GetAllTasksByUserIdOutput, error)
	CloseTask(ctx context.Context, input dto.CloseTaskInput) error
}

type TelegramUserService interface {
	GetOrCreate(ctx context.Context, input dto.CreateTelegramUserInput) (dto.CreateTelegramUserOutput, error)
}
