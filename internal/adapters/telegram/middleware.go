package telegram

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"strconv"
)

type Middleware struct {
	session Session
}

func NewMiddleware(session Session) *Middleware {
	return &Middleware{session: session}
}

func (m *Middleware) InitSession(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		sChatID := strconv.FormatInt(update.Message.Chat.ID, 10)
		err := m.session.InitSession(ctx, sChatID)
		if err != nil {
			log.Println("initSession error: ", err)
		}
		next(ctx, b, update)
	}
}
