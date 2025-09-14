package telegram

import (
	"github.com/go-telegram/bot"
	"github.com/shizakira/daily-tg-bot/internal/usecase"
)

type BotCommand = string

const (
	StartCommand      BotCommand = "start"
	TaskCreateCommand BotCommand = "task_create"
	TaskAllCommand    BotCommand = "task_all"
)

type Bot struct {
	bot          *bot.Bot
	stateMachine *BotTaskStateMachine
	service      *usecase.TaskService
}

func NewBot(bot *bot.Bot, stateMachine *BotTaskStateMachine, service *usecase.TaskService) *Bot {
	return &Bot{bot: bot, stateMachine: stateMachine, service: service}
}
