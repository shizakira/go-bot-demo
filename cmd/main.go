package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/shizakira/daily-tg-bot/config"
	"github.com/shizakira/daily-tg-bot/internal/adapters/postgres"
	"github.com/shizakira/daily-tg-bot/internal/adapters/redis"
	"github.com/shizakira/daily-tg-bot/internal/adapters/telegram"
	"github.com/shizakira/daily-tg-bot/internal/usecase"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func RunApp(c *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// logger
	logrus.SetOutput(os.Stdout)

	// database
	redisClient, err := redis.NewRedisClient(ctx, c.Redis)
	if err != nil {
		panic(err)
	}
	postgresPool, err := postgres.NewPostgresPool(c.Postgres)
	if err != nil {
		panic(err)
	}

	defer func(redisClient *redis.Redis) {
		_ = redisClient.Close()
	}(redisClient)
	defer func(postgresPool *postgres.Pool) {
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
	tgService := usecase.NewTelegramUserService(tgRepo, userRepo, taskRepo)

	// middleware
	tm := telegram.NewMiddleware(session, tgService)

	opts := []bot.Option{
		bot.WithMiddlewares(tm.GetMiddlewares()...),
		bot.WithWorkers(2),
	}
	b, err := bot.New(c.TgBot.Token, opts...)
	if err != nil {
		panic(err)
	}

	// handlers
	telegram.NewBot(b, session, taskService, tgService).InitHandlers()

	notifier := telegram.NewNotifier(b, tgRepo, taskRepo)

	go initScheduler(ctx, notifier)

	b.Start(ctx)

}

func main() {
	conf := config.Load()
	logrus.Info("starting tg bot")
	RunApp(conf)
}
