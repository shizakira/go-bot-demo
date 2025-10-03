package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/shizakira/daily-tg-bot/internal/adapters/postgres"
	"github.com/shizakira/daily-tg-bot/internal/adapters/redis"
)

type TelegramBotConfig struct {
	Token   string
	Workers int
}

type SchedulerConfig struct {
	Interval time.Duration
}

type Config struct {
	TgBot     TelegramBotConfig
	Redis     redis.Config
	Postgres  postgres.Config
	Scheduler SchedulerConfig
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	cfg := &Config{}
	var errs []error

	cfg.TgBot.Token = os.Getenv("TELEGRAM_BOT_TOKEN")
	if cfg.TgBot.Token == "" {
		errs = append(errs, errors.New("TELEGRAM_BOT_TOKEN is required"))
	}

	workersRaw := os.Getenv("TELEGRAM_BOT_WORKERS")
	if workersRaw != "" {
		workers, err := strconv.Atoi(workersRaw)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid TELEGRAM_BOT_WORKERS: %w", err))
		} else {
			cfg.TgBot.Workers = workers
		}
	}

	cfg.Redis = redis.Config{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
	if cfg.Redis.Addr == "" {
		errs = append(errs, errors.New("REDIS_ADDR is required"))
	}

	cfg.Postgres = postgres.Config{
		DSN: os.Getenv("DSN"),
	}
	if cfg.Postgres.DSN == "" {
		errs = append(errs, errors.New("DSN is required"))
	}

	intervalRaw := os.Getenv("SCHEDULER_INTERVAL")
	if intervalRaw == "" {
		cfg.Scheduler.Interval = 30 * time.Minute
	} else {
		dur, err := time.ParseDuration(intervalRaw)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid SCHEDULER_INTERVAL: %w", err))
		} else {
			cfg.Scheduler.Interval = dur
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return cfg, nil
}
