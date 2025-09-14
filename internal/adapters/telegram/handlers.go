package telegram

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
)

func (tb *Bot) onStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Printf("Message from %s with id %d", update.Message.From.Username, update.Message.From.ID)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Tap on /task-create to create new task",
	})
	if err != nil {
		log.Println("Error on sending message", err)
	}
}

func (tb *Bot) onGetTasks(ctx context.Context, b *bot.Bot, update *models.Update) {
	tasks, err := tb.service.GetAllTasks(ctx)
	if err != nil {
		log.Println("Error on sending message", err)
		return
	}
	for _, task := range tasks {
		msg := fmt.Sprintf(
			"ID: %d\nTitle: %s\nDescription: %s\nDeadline %s\n\n",
			task.ID, task.Title, task.Description, task.DeadlineDate.Format("2006-01-02 15:04"),
		)
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   msg,
		})
		if err != nil {
			log.Println("Error on sending message", err)
		}
	}
}

func (tb *Bot) onTaskCreate(ctx context.Context, b *bot.Bot, update *models.Update) {
	msg, err := tb.stateMachine.ProcessingCreateTask(ctx, update)
	if err != nil {
		log.Println("onTaskCreate", err)
		return
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
	if err != nil {
		log.Println("onTaskCreate", err)
	}
}
