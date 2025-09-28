package config

import (
	"github.com/joho/godotenv"
	"github.com/shizakira/daily-tg-bot/internal/adapters/postgres"
	"github.com/shizakira/daily-tg-bot/internal/adapters/redis"
	"os"
)

type TelegramBotConfig struct {
	Token string
}
type Config struct {
	TgBot    TelegramBotConfig
	Redis    redis.Config
	Postgres postgres.Config
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		TgBot: TelegramBotConfig{Token: os.Getenv("TELEGRAM_BOT_TOKEN")},
		Redis: redis.Config{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		Postgres: postgres.Config{
			DSN: os.Getenv("DSN"),
		},
	}
}
