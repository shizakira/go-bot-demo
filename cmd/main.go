package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/shizakira/daily-tg-bot/config"
	"github.com/shizakira/daily-tg-bot/internal/adapters/postgres"
	"github.com/shizakira/daily-tg-bot/internal/adapters/redis"
	"github.com/shizakira/daily-tg-bot/internal/adapters/telegram"
	"github.com/shizakira/daily-tg-bot/internal/usecase"
	"github.com/shizakira/daily-tg-bot/pkg/database"
	"log"
	"os"
	"os/signal"
)

func RunApp(c *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// database
	redisClient, err := database.NewRedisClient(ctx, c.Redis)
	if err != nil {
		panic(err)
	}
	postgresPool, err := database.NewPostgresPool(c.Postgres)
	if err != nil {
		panic(err)
	}

	defer func(redisClient *database.Redis) {
		_ = redisClient.Close()
	}(redisClient)
	defer func(postgresPool *database.PostgresPool) {
		_ = postgresPool.Close()
	}(postgresPool)

	// adapters
	session := redis.NewTelegramSession(redisClient)

	// repositories
	taskRepo := postgres.NewTaskRepository(postgresPool)
	tgRepo := postgres.NewTelegramUseRepository(postgresPool)
	userRepo := postgres.NewUserRepository(postgresPool)

	// services
	taskService := usecase.NewTaskService(taskRepo)
	tgService := usecase.NewTelegramUserService(tgRepo, userRepo)

	// telegram
	stmtM := telegram.NewBotTaskStateMachine(session, taskService)

	// middleware
	tm := telegram.NewMiddleware(session, tgService)

	opts := []bot.Option{bot.WithMiddlewares(tm.InitTGUserSession)}
	b, err := bot.New(c.TgBot.Token, opts...)
	if err != nil {
		panic(err)
	}

	// handlers
	telegram.NewBot(b, stmtM, taskService).InitHandlers()

	b.Start(ctx)
}

func main() {
	conf := config.Load()
	log.Println("Starting telegram bot")
	RunApp(conf)
}
