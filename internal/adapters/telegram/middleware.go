package telegram

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/shizakira/daily-tg-bot/internal/dto"
	"github.com/shizakira/daily-tg-bot/internal/usecase"
	"log"
	"strconv"
)

type Middleware struct {
	session   Session
	tgService *usecase.TelegramUserService
}

type UserIdIdempotencyKey string

func NewMiddleware(session Session, userService *usecase.TelegramUserService) *Middleware {
	return &Middleware{
		session:   session,
		tgService: userService,
	}
}

func (m *Middleware) InitTGUserSession(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var chatID int64
		var user *models.User

		if update.Message != nil {
			chatID = update.Message.Chat.ID
			user = update.Message.From
		} else if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Message.Chat.ID
			user = &update.CallbackQuery.From
		} else {
			next(ctx, b, update)
			return
		}

		sChatID := strconv.FormatInt(chatID, 10)
		if err := m.session.InitSession(ctx, sChatID); err != nil {
			log.Println("initSession error:", err)
		}

		output, err := m.tgService.GetOrCreate(ctx, dto.CreateTelegramUserInput{
			ChatID:     chatID,
			TelegramID: user.ID,
			Username:   user.Username,
		})
		if err != nil {
			log.Println("GetOrCreate error:", err)
		}

		newCtx := context.WithValue(ctx, UserIdIdempotencyKey("userId"), output.UserID)
		log.Println("InitTGUserSession userId:", newCtx.Value(UserIdIdempotencyKey("userId")))
		log.Println("test we are here")
		next(newCtx, b, update)
	}
}
