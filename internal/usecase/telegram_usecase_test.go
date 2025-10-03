package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/internal/dto"
)

type telegramUserRepoMock struct {
	findResult *domain.TelegramUser
	findErr    error
	created    []*domain.TelegramUser
	createErr  error
}

func (m *telegramUserRepoMock) Create(_ context.Context, user *domain.TelegramUser) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.created = append(m.created, user)
	return nil
}

func (m *telegramUserRepoMock) FindByChatID(_ context.Context, _ int64) (*domain.TelegramUser, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.findResult, nil
}

func (m *telegramUserRepoMock) FindByUserIDs(context.Context, []int64) ([]*domain.TelegramUser, error) {
	return nil, nil
}

type userRepoMock struct {
	created   []*domain.User
	createErr error
	nextID    int64
}

func (m *userRepoMock) Create(_ context.Context, user *domain.User) (int64, error) {
	m.created = append(m.created, user)
	if m.createErr != nil {
		return 0, m.createErr
	}
	return m.nextID, nil
}

func (m *userRepoMock) FindByID(context.Context, int64) (*domain.User, error) {
	return nil, errors.New("not implemented")
}

func (m *userRepoMock) Exists(context.Context, int64) (bool, error) {
	return false, errors.New("not implemented")
}

func TestTelegramUserService_GetOrCreate_UserExists(t *testing.T) {
	tgRepo := &telegramUserRepoMock{findResult: &domain.TelegramUser{UserID: 77}}
	userRepo := &userRepoMock{nextID: 88}
	service := NewTelegramUserService(tgRepo, userRepo)

	output, err := service.GetOrCreate(context.Background(), dto.CreateTelegramUserInput{ChatID: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.UserID != 77 {
		t.Fatalf("expected user id 77, got %d", output.UserID)
	}
	if len(userRepo.created) != 0 {
		t.Fatalf("expected no new users created, got %d", len(userRepo.created))
	}
	if len(tgRepo.created) != 0 {
		t.Fatalf("expected telegram user not created, got %d", len(tgRepo.created))
	}
}

func TestTelegramUserService_GetOrCreate_CreateNewUser(t *testing.T) {
	tgRepo := &telegramUserRepoMock{findErr: sql.ErrNoRows}
	userRepo := &userRepoMock{nextID: 101}
	service := NewTelegramUserService(tgRepo, userRepo)

	input := dto.CreateTelegramUserInput{ChatID: 1, TelegramID: 2, Username: "john"}
	output, err := service.GetOrCreate(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(userRepo.created) != 1 {
		t.Fatalf("expected 1 user created, got %d", len(userRepo.created))
	}
	if len(tgRepo.created) != 1 {
		t.Fatalf("expected 1 telegram user created, got %d", len(tgRepo.created))
	}

	created := tgRepo.created[0]
	if created.UserID != userRepo.nextID {
		t.Errorf("expected telegram user userID %d, got %d", userRepo.nextID, created.UserID)
	}
	if created.ChatID != input.ChatID {
		t.Errorf("expected chat id %d, got %d", input.ChatID, created.ChatID)
	}
	if created.TelegramID != input.TelegramID {
		t.Errorf("expected telegram id %d, got %d", input.TelegramID, created.TelegramID)
	}
	if created.Username != input.Username {
		t.Errorf("expected username %q, got %q", input.Username, created.Username)
	}

	if output.UserID != userRepo.nextID {
		t.Fatalf("expected output user id %d, got %d", userRepo.nextID, output.UserID)
	}
}
