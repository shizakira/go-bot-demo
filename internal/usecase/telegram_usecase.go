package usecase

import (
	"context"
	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/internal/dto"
	"github.com/shizakira/daily-tg-bot/internal/ports"
)

type TelegramUserService struct {
	telegramRepo ports.TelegramUserRepository
	userRepo     ports.UserRepository
	taskRepo     ports.TaskRepository
}

func NewTelegramUserService(
	telegramRepo ports.TelegramUserRepository,
	userRepo ports.UserRepository,
	taskRepo ports.TaskRepository,
) *TelegramUserService {
	return &TelegramUserService{telegramRepo: telegramRepo, userRepo: userRepo, taskRepo: taskRepo}
}

func (tu *TelegramUserService) GetOrCreate(
	ctx context.Context,
	input dto.CreateTelegramUserInput,
) (dto.CreateTelegramUserOutput, error) {
	tgUser, err := tu.telegramRepo.FindByChatID(ctx, input.ChatID)
	output := dto.CreateTelegramUserOutput{}
	if err == nil {
		output.UserID = tgUser.UserID
		return output, nil
	}

	userID, err := tu.userRepo.Create(ctx, new(domain.User))
	if err != nil {
		return output, err
	}
	tgUser = &domain.TelegramUser{
		UserID:     userID,
		ChatID:     input.ChatID,
		TelegramID: input.TelegramID,
		Username:   input.Username,
	}
	err = tu.telegramRepo.Create(ctx, tgUser)
	if err != nil {
		return output, err
	}
	return output, nil
}
