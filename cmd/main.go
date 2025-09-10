package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/shizakira/daily-tg-bot/config"
	"github.com/shizakira/daily-tg-bot/internal/adapters"
	"github.com/shizakira/daily-tg-bot/internal/telegram"
	"github.com/shizakira/daily-tg-bot/pkg/database"
	"log"
	"os"
	"os/signal"
)

func RunApp(c *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	//opts := []bot.Option{}
	b, err := bot.New(c.TgBot.Token)
	if err != nil {
		panic(err)
	}

	// database
	redisClient, err := database.NewRedisClient(ctx, c.Redis)
	if err != nil {
		panic(err)
	}

	defer redisClient.Close()

	// adapters
	redisStore := adapters.NewRedisStore[telegram.TaskState](redisClient)

	// action
	actions := telegram.NewAction(redisStore)

	// handlers
	telegram.NewBotHook(actions).InitBotHandlers(b)

	b.Start(ctx)
}

func main() {
	conf := config.Load()
	log.Println("Starting telegram bot")
	RunApp(conf)
}
