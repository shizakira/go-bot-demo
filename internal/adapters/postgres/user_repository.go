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
	stmt, err := ur.pool.PrepareContext(ctx, "insert into users default values returning id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var userID int64
	err = stmt.QueryRowContext(ctx).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (ur *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	stmt, err := ur.pool.PrepareContext(ctx, "select id from users where id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	user := &domain.User{}
	err = stmt.QueryRowContext(ctx, id).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) Exists(ctx context.Context, id int64) (bool, error) {
	stmt, err := ur.pool.PrepareContext(ctx, "select exists(select 1 from users where id = $1)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var exists bool
	err = stmt.QueryRowContext(ctx, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
