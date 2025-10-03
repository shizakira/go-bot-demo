package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/shizakira/daily-tg-bot/config"
	"github.com/shizakira/daily-tg-bot/internal/adapters/postgres"
	"github.com/shizakira/daily-tg-bot/internal/adapters/redis"
	tg "github.com/shizakira/daily-tg-bot/internal/adapters/telegram"
	"github.com/shizakira/daily-tg-bot/internal/usecase"
	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg       *config.Config
	redis     *redis.Redis
	postgres  *postgres.Pool
	bot       *bot.Bot
	scheduler *NotifierScheduler
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	redisClient, err := redis.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("init redis client: %w", err)
	}

	postgresPool, err := postgres.NewPostgresPool(cfg.Postgres)
	if err != nil {
		_ = redisClient.Close()
		return nil, fmt.Errorf("init postgres pool: %w", err)
	}

	session := redis.NewTelegramSession(redisClient)

	taskRepo := postgres.NewTaskRepository(postgresPool)
	tgRepo := postgres.NewTelegramUseRepository(postgresPool)
	userRepo := postgres.NewUserRepository(postgresPool)

	taskService := usecase.NewTaskService(taskRepo)
	tgService := usecase.NewTelegramUserService(tgRepo, userRepo)

	middleware := tg.NewMiddleware(session, tgService)

	workers := cfg.TgBot.Workers

	botOptions := []bot.Option{
		bot.WithMiddlewares(middleware.GetMiddlewares()...),
		bot.WithWorkers(workers),
	}

	telegramBot, err := bot.New(cfg.TgBot.Token, botOptions...)
	if err != nil {
		_ = redisClient.Close()
		_ = postgresPool.Close()
		return nil, fmt.Errorf("init telegram bot: %w", err)
	}

	tg.NewBot(telegramBot, session, taskService, tgService).InitHandlers()

	notifier := tg.NewNotifier(telegramBot, tgRepo, taskRepo)
	sched, err := NewNotifierScheduler(notifier, cfg.Scheduler.Interval)
	if err != nil {
		_ = redisClient.Close()
		_ = postgresPool.Close()
		return nil, fmt.Errorf("init scheduler: %w", err)
	}

	return &App{
		cfg:       cfg,
		redis:     redisClient,
		postgres:  postgresPool,
		bot:       telegramBot,
		scheduler: sched,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return a.scheduler.Start(ctx)
	})

	group.Go(func() error {
		a.bot.Start(ctx)
		return ctx.Err()
	})

	return group.Wait()
}

func (a *App) Close() error {
	var errs []error

	if a.redis != nil {
		if err := a.redis.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close redis: %w", err))
		}
	}

	if a.postgres != nil {
		if err := a.postgres.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close postgres: %w", err))
		}
	}

	return errors.Join(errs...)
}
