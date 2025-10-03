package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/internal/dto"
	"github.com/shizakira/daily-tg-bot/internal/ports"
)

type TelegramUserService struct {
	telegramRepo ports.TelegramUserRepository
	userRepo     ports.UserRepository
}

func NewTelegramUserService(
	telegramRepo ports.TelegramUserRepository,
	userRepo ports.UserRepository,
) *TelegramUserService {
	return &TelegramUserService{telegramRepo: telegramRepo, userRepo: userRepo}
}

func (tu *TelegramUserService) GetOrCreate(
	ctx context.Context,
	input dto.CreateTelegramUserInput,
) (dto.CreateTelegramUserOutput, error) {
	output := dto.CreateTelegramUserOutput{}

	tgUser, err := tu.telegramRepo.FindByChatID(ctx, input.ChatID)
	if err == nil {
		output.UserID = tgUser.UserID
		return output, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return output, err
	}

	userID, err := tu.userRepo.Create(ctx, new(domain.User))
	if err != nil {
		return output, err
	}

	newTelegramUser := &domain.TelegramUser{
		UserID:     userID,
		ChatID:     input.ChatID,
		TelegramID: input.TelegramID,
		Username:   input.Username,
	}
	if err := tu.telegramRepo.Create(ctx, newTelegramUser); err != nil {
		return output, err
	}

	output.UserID = userID
	return output, nil
}
