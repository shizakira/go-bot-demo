package config

import (
	"github.com/joho/godotenv"
	"github.com/shizakira/daily-tg-bot/pkg/database"
	"log"
	"os"
)

type TelegramBotConfig struct {
	Token string
}
type Config struct {
	TgBot    TelegramBotConfig
	Redis    database.RedisConfig
	Postgres database.PostgresConfig
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		TgBot: TelegramBotConfig{Token: os.Getenv("TELEGRAM_BOT_TOKEN")},
		Redis: database.RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		Postgres: database.PostgresConfig{
			DSN: os.Getenv("DSN"),
		},
	}
}
