package usecase

import (
	"context"
	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/internal/dto"
)

type TelegramUserRepository interface {
	Create(ctx context.Context, user *domain.TelegramUser) error
	FindByChatID(ctx context.Context, chatId int64) (*domain.TelegramUser, error)
}

type TelegramUserService struct {
	tgRepo TelegramUserRepository
	uRepo  UserRepository
}

func NewTelegramUserService(tgRepo TelegramUserRepository, uRepo UserRepository) *TelegramUserService {
	return &TelegramUserService{tgRepo: tgRepo, uRepo: uRepo}
}

func (tu *TelegramUserService) GetOrCreate(ctx context.Context, input dto.CreateTelegramUser) (int64, error) {
	tgUser, err := tu.tgRepo.FindByChatID(ctx, input.ChatID)
	if err == nil {
		return tgUser.UserID, nil
	}

	userID, err := tu.uRepo.Create(ctx, new(domain.User))
	if err != nil {
		return 0, err
	}
	tgUser = &domain.TelegramUser{
		UserID:     userID,
		ChatID:     input.ChatID,
		TelegramID: input.TelegramID,
		Username:   input.Username,
	}
	err = tu.tgRepo.Create(ctx, tgUser)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
