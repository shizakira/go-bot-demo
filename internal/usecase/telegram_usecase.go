package usecase

import (
	"context"
	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/internal/dto"
	"github.com/shizakira/daily-tg-bot/internal/ports"
)

type TelegramUserService struct {
	tgRepo ports.TelegramUserRepository
	uRepo  ports.UserRepository
}

func NewTelegramUserService(tgRepo ports.TelegramUserRepository, uRepo ports.UserRepository) *TelegramUserService {
	return &TelegramUserService{tgRepo: tgRepo, uRepo: uRepo}
}

func (tu *TelegramUserService) GetOrCreate(ctx context.Context, input dto.CreateTelegramUserInput) (dto.CreateTelegramUserOutput, error) {
	tgUser, err := tu.tgRepo.FindByChatID(ctx, input.ChatID)
	output := dto.CreateTelegramUserOutput{}
	if err == nil {
		output.UserID = tgUser.UserID
		return output, nil
	}

	userID, err := tu.uRepo.Create(ctx, new(domain.User))
	if err != nil {
		return output, err
	}
	tgUser = &domain.TelegramUser{
		UserID:     userID,
		ChatID:     input.ChatID,
		TelegramID: input.TelegramID,
		Username:   input.Username,
	}
	err = tu.tgRepo.Create(ctx, tgUser)
	if err != nil {
		return output, err
	}
	return output, nil
}
