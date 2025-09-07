package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type TelegramBotConfig struct {
	Token string
}
type Config struct {
	TgBot TelegramBotConfig
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tgConf := TelegramBotConfig{Token: os.Getenv("TELEGRAM_BOT_TOKEN")}

	return &Config{
		TgBot: tgConf,
	}
}
