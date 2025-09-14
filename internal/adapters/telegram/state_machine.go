package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-telegram/bot/models"
	"github.com/shizakira/daily-tg-bot/internal/dto"
	"github.com/shizakira/daily-tg-bot/internal/usecase"
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
	NextStep step                `json:"next_step"`
	Data     dto.CreateTaskInput `json:"data"`
}

type BotTaskStateMachine struct {
	session Session
	service *usecase.TaskService
}

func NewBotTaskStateMachine(session Session, service *usecase.TaskService) *BotTaskStateMachine {
	return &BotTaskStateMachine{session: session, service: service}
}

func (st *BotTaskStateMachine) ProcessingCreateTask(ctx context.Context, update *models.Update) (string, error) {
	sChatID := strconv.FormatInt(update.Message.Chat.ID, 10)
	msgText := update.Message.Text
	userId, ok := ctx.Value(UserIdIdempotencyKey("userId")).(int64)
	if !ok {
		return "", errors.New("something wrong with ctx user id")
	}

	if update.Message.Text == "/"+taskCreateCommand {
		state := TaskState{
			NextStep: titleStep,
			Data:     dto.CreateTaskInput{UserID: userId},
		}
		err := st.session.Set(ctx, sChatID, "task_creating_state", &state)
		if err != nil {
			return "", err
		}
		return "Enter the task title", nil
	}
	stateRaw, err := st.session.Get(ctx, sChatID, "task_creating_state")
	if err != nil {
		return "", err
	}
	state := new(TaskState)
	err = json.Unmarshal(stateRaw, state)
	if err != nil {
		return "", err
	}

	switch state.NextStep {
	case titleStep:
		state.Data.Title = msgText
		state.NextStep = descStep
		err = st.session.Set(ctx, sChatID, "task_creating_state", state)
		if err != nil {
			return "", err
		}
		return "Enter the task description", nil
	case descStep:
		state.Data.Description = msgText
		state.NextStep = dateTimeStep
		err = st.session.Set(ctx, sChatID, "task_creating_state", state)
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
		err = st.session.Del(ctx, sChatID, "task_creating_state")
		if err != nil {
			return "", err
		}
		return "Task was created", nil
	default:
		return "", errors.New("not create task flow input")
	}
}
