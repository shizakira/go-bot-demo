package telegram

import (
	"github.com/go-telegram/bot"
	"github.com/shizakira/daily-tg-bot/internal/usecase"
)

type botCommand = string

const (
	startCommand      botCommand = "start"
	taskCreateCommand botCommand = "task_create"
	taskAllCommand    botCommand = "task_all"
)

type Bot struct {
	bot         *bot.Bot
	session     Session
	taskService *usecase.TaskService
	tgService   *usecase.TelegramUserService
}

func NewBot(bot *bot.Bot, session Session, taskService *usecase.TaskService, tgService *usecase.TelegramUserService) *Bot {
	return &Bot{bot: bot, session: session, taskService: taskService, tgService: tgService}
}

func (tb *Bot) InitHandlers() {
	tb.bot.RegisterHandler(bot.HandlerTypeMessageText, startCommand, bot.MatchTypeCommandStartOnly, tb.onStart)
	tb.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, taskCreateCommand, bot.MatchTypeCommand, tb.onTaskCreate)
	tb.bot.RegisterHandler(bot.HandlerTypeMessageText, taskAllCommand, bot.MatchTypeCommand, tb.onGetTasks)
	tb.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, "button_cancel", bot.MatchTypeExact, tb.onTaskCancel)
	tb.bot.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, tb.onTaskCreate)
}
