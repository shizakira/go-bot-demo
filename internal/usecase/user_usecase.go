package usecase

import (
	"context"
	"github.com/shizakira/daily-tg-bot/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (int64, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	Exists(ctx context.Context, id int64) (bool, error)
}
