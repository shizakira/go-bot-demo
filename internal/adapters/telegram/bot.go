package telegram

import (
	"github.com/go-telegram/bot"
	"github.com/shizakira/daily-tg-bot/internal/ports"
)

type botCommand = string

const (
	startCommand      botCommand = "start"
	taskCreateCommand botCommand = "task_create"
	taskAllCommand    botCommand = "task_all"
)

type botButton = string

const (
	taskCancelButton botButton = "task_cancel_btn"
	taskDone         botButton = "task_done_btn"
	taskClose        botButton = "task_close_btn"
)

type Bot struct {
	bot         *bot.Bot
	session     Session
	taskService ports.TaskService
	tgService   ports.TelegramUserService
}

func NewBot(bot *bot.Bot, session Session, taskService ports.TaskService, tgService ports.TelegramUserService) *Bot {
	return &Bot{bot: bot, session: session, taskService: taskService, tgService: tgService}
}

func (tb *Bot) InitHandlers() {
	tb.bot.RegisterHandler(bot.HandlerTypeMessageText, startCommand, bot.MatchTypeCommandStartOnly, tb.onStart)
	tb.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, taskCreateCommand, bot.MatchTypeCommand, tb.onTaskCreate)
	tb.bot.RegisterHandler(bot.HandlerTypeMessageText, taskAllCommand, bot.MatchTypeCommand, tb.onGetTasks)
	tb.bot.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, tb.onTaskCreate)

	tb.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, taskCancelButton, bot.MatchTypeExact, tb.onTaskCancel)
	tb.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, taskDone, bot.MatchTypeContains, tb.onTaskDone)
	tb.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, taskClose, bot.MatchTypeContains, tb.onTaskClose)
}
