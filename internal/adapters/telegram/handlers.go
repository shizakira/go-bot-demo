package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/shizakira/daily-tg-bot/internal/dto"
	"log"
	"strconv"
	"time"
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
	userId, ok := ctx.Value(UserIdIdempotencyKey("userId")).(int64)
	if !ok {
		log.Println("Error on sending message", fmt.Errorf("something wrong with ctx user id: %d", userId))
		return
	}
	output, err := tb.taskService.GetAllTasksByUserID(ctx, dto.GetAllTasksByUserIdInput{UserID: userId})
	if err != nil {
		log.Println("Error on sending message", err)
		return
	}
	for _, task := range output.Tasks {
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

func (tb *Bot) onTaskCancel(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	sChatID := strconv.FormatInt(chatID, 10)

	if err := tb.session.Del(ctx, sChatID, "task_creating_state"); err != nil {
		log.Println("error deleting state:", err)
	}

	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
	if err != nil {
		log.Println("onTaskCancel", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "❌ Создание задачи отменено",
	})
	if err != nil {
		log.Println("onTaskCancel", err)
	}
}

func (tb *Bot) onTaskCreate(ctx context.Context, b *bot.Bot, update *models.Update) {
	if err := tb.processingCreateTask(ctx, b, update); err != nil {
		log.Println("onTaskCreate", err)
	}
}

func (tb *Bot) getCancelCreatingTaskKB() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Cancel", CallbackData: "button_cancel"},
			},
		},
	}
}

func (tb *Bot) processingCreateTask(ctx context.Context, b *bot.Bot, update *models.Update) error {
	sChatID := strconv.FormatInt(update.Message.Chat.ID, 10)
	msgText := update.Message.Text
	userId, ok := ctx.Value(UserIdIdempotencyKey("userId")).(int64)
	if !ok {
		return errors.New("something wrong with ctx user id")
	}

	if msgText == "/"+taskCreateCommand {
		state := TaskState{
			NextStep: titleStep,
			Data:     dto.CreateTaskInput{UserID: userId},
		}
		if err := tb.session.Set(ctx, sChatID, "task_creating_state", &state); err != nil {
			return err
		}
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Enter the task title",
			ReplyMarkup: tb.getCancelCreatingTaskKB(),
		})
		return err
	}

	stateRaw, err := tb.session.Get(ctx, sChatID, "task_creating_state")
	if err != nil {
		return err
	}

	state := new(TaskState)
	if err := json.Unmarshal(stateRaw, state); err != nil {
		return err
	}

	switch state.NextStep {
	case titleStep:
		if vErr := validateTitle(msgText); vErr != nil {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        vErr.Error(),
				ReplyMarkup: tb.getCancelCreatingTaskKB(),
			})
			return err
		}

		state.Data.Title = msgText
		state.NextStep = descStep
		if err := tb.session.Set(ctx, sChatID, "task_creating_state", state); err != nil {
			return err
		}

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Enter the task description",
			ReplyMarkup: tb.getCancelCreatingTaskKB(),
		})
		return err

	case descStep:
		if vErr := validateDesc(msgText); vErr != nil {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        vErr.Error(),
				ReplyMarkup: tb.getCancelCreatingTaskKB(),
			})
			return err
		}

		state.Data.Description = msgText
		state.NextStep = dateTimeStep
		if err := tb.session.Set(ctx, sChatID, "task_creating_state", state); err != nil {
			return err
		}

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Enter the task deadline datetime",
			ReplyMarkup: tb.getCancelCreatingTaskKB(),
		})
		return err

	case dateTimeStep:
		loc, _ := time.LoadLocation("Asia/Yekaterinburg")
		t1, err := time.ParseInLocation("2006-01-02 15:04", msgText, loc)
		if vErr := validateDate(t1); vErr != nil {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        vErr.Error(),
				ReplyMarkup: tb.getCancelCreatingTaskKB(),
			})
			return err
		}

		state.Data.DeadlineDate = t1
		if err := tb.taskService.CreateTask(ctx, state.Data); err != nil {
			return err
		}

		if err := tb.session.Del(ctx, sChatID, "task_creating_state"); err != nil {
			return err
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Task was created",
		})
		return err

	default:
		return errors.New("not create task flow input")
	}
}
