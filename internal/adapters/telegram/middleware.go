package telegram

import (
	"context"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/shizakira/daily-tg-bot/internal/dto"
	"github.com/shizakira/daily-tg-bot/internal/ports"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	session   Session
	tgService ports.TelegramUserService
}

type UserIdIdempotencyKey string

func NewMiddleware(session Session, userService ports.TelegramUserService) *Middleware {
	return &Middleware{
		session:   session,
		tgService: userService,
	}
}

func (m *Middleware) GetMiddlewares() []bot.Middleware {
	return []bot.Middleware{m.RequestMiddleware, m.InitTGUserSession}
}

func (m *Middleware) RequestMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		logrus.WithFields(logrus.Fields{
			"telegram_user_id": getUser(update).ID,
			"chat_id":          getChat(update).ID,
		}).Info("requesting user")
		next(ctx, b, update)
	}
}

func (m *Middleware) InitTGUserSession(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chat := getChat(update)
		user := getUser(update)

		if err := m.session.InitSession(ctx, strconv.FormatInt(chat.ID, 10)); err != nil {
			logrus.Error("initSession error:", err)
		}

		output, err := m.tgService.GetOrCreate(ctx, dto.CreateTelegramUserInput{
			ChatID:     chat.ID,
			TelegramID: user.ID,
			Username:   user.Username,
		})
		if err != nil {
			logrus.Error("GetOrCreate error:", err)
		}

		newCtx := context.WithValue(ctx, UserIdIdempotencyKey("userId"), output.UserID)
		next(newCtx, b, update)
	}
}
