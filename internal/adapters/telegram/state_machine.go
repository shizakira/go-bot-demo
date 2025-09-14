package telegram

import (
	"context"
	"errors"
	"github.com/go-telegram/bot/models"
	"github.com/shizakira/daily-tg-bot/internal/dto"
	"github.com/shizakira/daily-tg-bot/internal/usecase"
	"log"
	"strconv"
	"time"
)

type step string

const (
	titleStep    step = "title"
	descStep     step = "description"
	dateTimeStep step = "datetime"
)

type TaskState struct {
	UserID   int64               `json:"user_id"`
	NextStep step                `json:"next_step"`
	Data     dto.CreateTaskInput `json:"data"`
}

type StateStore interface {
	Set(ctx context.Context, key string, value *TaskState) error
	Get(ctx context.Context, key string) (*TaskState, error)
	Del(ctx context.Context, key string)
}

type BotTaskStateMachine struct {
	store   StateStore
	service *usecase.TaskService
}

func NewBotTaskStateMachine(store StateStore, service *usecase.TaskService) *BotTaskStateMachine {
	return &BotTaskStateMachine{store: store, service: service}
}

func (st *BotTaskStateMachine) ProcessingCreateTask(ctx context.Context, update *models.Update) (string, error) {
	sChatID := strconv.FormatInt(update.Message.Chat.ID, 10)
	msgText := update.Message.Text
	log.Println(update.Message.Text)
	if update.Message.Text == "/"+taskCreateCommand {
		state := TaskState{
			UserID:   update.Message.From.ID,
			NextStep: titleStep,
			Data:     dto.CreateTaskInput{},
		}
		err := st.store.Set(ctx, sChatID, &state)
		if err != nil {
			return "", err
		}
		return "Enter the task title", nil
	}
	state, err := st.store.Get(ctx, sChatID)
	if err != nil {
		return "", err
	}

	switch state.NextStep {
	case titleStep:
		state.Data.Title = msgText
		state.NextStep = descStep
		err = st.store.Set(ctx, sChatID, state)
		if err != nil {
			return "", err
		}
		return "Enter the task description", nil
	case descStep:
		state.Data.Description = msgText
		state.NextStep = dateTimeStep
		err = st.store.Set(ctx, sChatID, state)
		if err != nil {
			return "", err
		}
		return "Enter the task deadline date in format 2025-09-01 00:00", nil
	case dateTimeStep:
		parsedDate, err := time.Parse("2006-01-02 15:04", msgText)
		if err != nil {
			return "", err
		}
		state.Data.DeadlineDate = parsedDate
		err = st.service.CreateTask(ctx, state.Data)
		if err != nil {
			return "", err
		}
		st.store.Del(ctx, sChatID)
		return "Task was created", nil
	default:
		return "", errors.New("not create task flow input")
	}
}
