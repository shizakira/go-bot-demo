package telegram

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/datepicker"
	"log"
	"time"
)

type BotHook struct {
	action *Action
}

func NewBotHook(action *Action) *BotHook {
	return &BotHook{action: action}
}

func (h *BotHook) onStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Printf("Message from %s with id %d", update.Message.From.Username, update.Message.From.ID)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Tap on /task-create to create new task",
	})
	if err != nil {
		log.Println("Error on sending message", err)
	}
}

func (h *BotHook) onTaskCreate(ctx context.Context, b *bot.Bot, update *models.Update) {
	msgParams, err := h.action.CreateTaskFlow(ctx, &TaskDTO{
		ChatID: update.Message.Chat.ID,
		UserID: update.Message.From.ID,
		Init:   true,
	})
	if err != nil {
		log.Println("Error on sending message", err)
		return
	}
	_, err = b.SendMessage(ctx, msgParams)
	if err != nil {
		log.Println("Error on sending message", err)
	}
}

func (h *BotHook) onTextInput(ctx context.Context, b *bot.Bot, update *models.Update) {
	msgParams, err := h.action.CreateTaskFlow(ctx, &TaskDTO{
		ChatID:      update.Message.Chat.ID,
		MsgText:     update.Message.Text,
		ReplyMarkup: datepicker.New(b, h.onDatepickerSimpleSelect),
	})
	if err != nil {
		log.Println("Error on sending message", err)
		return
	}
	_, err = b.SendMessage(ctx, msgParams)
	if err != nil {
		log.Println("Error on sending message", err)
	}
}

func (h *BotHook) onDatepickerSimpleSelect(
	ctx context.Context,
	b *bot.Bot,
	mes models.MaybeInaccessibleMessage,
	date time.Time,
) {
	msgParams, err := h.action.CreateTaskFlow(ctx, &TaskDTO{
		ChatID: mes.Message.Chat.ID,
		Date:   date,
	})
	if err != nil {
		log.Println("Error on sending message", err)
		return
	}
	_, err = b.SendMessage(ctx, msgParams)
	if err != nil {
		log.Println("Error on sending message", err)
	}
}

func (h *BotHook) InitBotHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommandStartOnly, h.onStart)
	b.RegisterHandler(bot.HandlerTypeMessageText, "task_create", bot.MatchTypeCommand, h.onTaskCreate)
	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, h.onTextInput)
}
