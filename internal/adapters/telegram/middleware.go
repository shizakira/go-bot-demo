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
		sChatID := strconv.FormatInt(update.Message.Chat.ID, 10)
		err := m.session.InitSession(ctx, sChatID)
		if err != nil {
			log.Println("initSession error: ", err)
		}
		output, err := m.tgService.GetOrCreate(ctx, dto.CreateTelegramUserInput{
			ChatID:     update.Message.Chat.ID,
			TelegramID: update.Message.From.ID,
			Username:   update.Message.From.Username,
		})
		if err != nil {
			log.Println("initSession error: ", err)
		}
		newCtx := context.WithValue(ctx, UserIdIdempotencyKey("userId"), output.UserID)
		log.Println(newCtx.Value(UserIdIdempotencyKey("userId")))
		next(newCtx, b, update)
	}
}
