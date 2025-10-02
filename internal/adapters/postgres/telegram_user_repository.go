package postgres

import (
	"context"
	"fmt"
	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/pkg/helpers"
)

type TelegramUseRepository struct {
	pool *Pool
}

func NewTelegramUseRepository(pool *Pool) *TelegramUseRepository {
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

func (tr *TelegramUseRepository) FindByUserIDs(ctx context.Context, userIDs []int64) ([]*domain.TelegramUser, error) {
	query := fmt.Sprintf(
		"select * from telegram_users where user_id in (%s);",
		helpers.GeneratePlaceholders(len(userIDs)),
	)
	args := make([]any, len(userIDs))
	for i, id := range userIDs {
		args[i] = id
	}
	//rows, err := tr.pool.QueryContext(ctx, query, args...)
	rows, err := tr.pool.Query(query, args...)

	var telegramUsers []*domain.TelegramUser
	for rows.Next() {
		tgUser := new(domain.TelegramUser)
		if err = rows.Scan(
			&tgUser.ID, &tgUser.UserID, &tgUser.TelegramID, &tgUser.ChatID, &tgUser.Username,
		); err != nil {
			return nil, err
		}
		telegramUsers = append(telegramUsers, tgUser)
	}

	return telegramUsers, nil
}
