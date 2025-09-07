package telegram

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
)

type BotHandler struct{}

func NewBotHandlers() *BotHandler {
	return &BotHandler{}
}

func (h *BotHandler) onStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Printf("Message from %s with id %d", update.Message.From.Username, update.Message.From.ID)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Дарова",
	})
	if err != nil {
		log.Println("Error on sending message", err)
		return
	}
}

func (h *BotHandler) onTaskCreate(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Введите название задачи",
	})
	if err != nil {
		log.Println("Error on sending message", err)
		return
	}
}

func (h *BotHandler) onTextInput(ctx context.Context, b *bot.Bot, update *models.Update) {
}

func (h *BotHandler) Init(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommandStartOnly, h.onStart)
	b.RegisterHandler(bot.HandlerTypeMessageText, "task_create", bot.MatchTypeCommand, h.onTaskCreate)
	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, h.onTextInput)
}
