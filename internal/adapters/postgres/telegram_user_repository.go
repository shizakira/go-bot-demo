package postgres

import (
	"context"
	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/pkg/database"
)

type TelegramUseRepository struct {
	pool *database.PostgresPool
}

func NewTelegramUseRepository(pool *database.PostgresPool) *TelegramUseRepository {
	return &TelegramUseRepository{pool: pool}
}

func (tr *TelegramUseRepository) Create(ctx context.Context, user *domain.TelegramUser) error {
	stmt, err := tr.pool.PrepareContext(ctx,
		"insert into telegram_users (user_id, chat_id, telegram_id, username) values ($1, $2, $3, $4)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, user.UserID, user.ChatID, user.TelegramID, user.Username).Scan(&user.ID)
	return err
}

func (tr *TelegramUseRepository) FindByChatID(ctx context.Context, chatId int64) (*domain.TelegramUser, error) {
	stmt, err := tr.pool.PrepareContext(ctx, "select * from telegram_users where chat_id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	tgUser := new(domain.TelegramUser)
	err = stmt.QueryRowContext(ctx, chatId).Scan(
		&tgUser.ID, &tgUser.UserID, &tgUser.TelegramID, &tgUser.ChatID, &tgUser.Username)
	if err != nil {
		return nil, err
	}
	return tgUser, nil
}
