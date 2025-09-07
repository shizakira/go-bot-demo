package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/shizakira/daily-tg-bot/config"
	"log"
	"os"
	"os/signal"
)

func RunApp(conf *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(conf.TgBot.Token, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func main() {
	conf := config.Load()
	log.Println("Starting telegram bot")
	RunApp(conf)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Printf("Message from %s with id %d", update.Message.From.FirstName, update.Message.From.ID)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
	if err != nil {
		log.Println("Error on sending message", err)
		return
	}
}
