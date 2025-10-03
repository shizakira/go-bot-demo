package postgres

import (
	"context"

	"github.com/shizakira/daily-tg-bot/internal/domain"
)

type UserRepository struct {
	pool *Pool
}

func NewUserRepository(pool *Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (ur *UserRepository) Create(ctx context.Context, user *domain.User) (int64, error) {
	var userID int64
	err := ur.pool.QueryRowContext(ctx, "insert into users default values returning id").Scan(&userID)
	return userID, err
}

func (ur *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	user := &domain.User{}
	if err := ur.pool.QueryRowContext(ctx, "select id from users where id = $1", id).Scan(&user.ID); err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := ur.pool.QueryRowContext(ctx, "select exists(select 1 from users where id = $1)", id).Scan(&exists)
	return exists, err
}
