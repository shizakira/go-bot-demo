package telegram

import (
	"context"
	"errors"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/datepicker"
	"strconv"
	"time"
)

type TaskDTO struct {
	ChatID      int64
	UserID      int64
	MsgText     string
	Init        bool
	Date        time.Time
	ReplyMarkup *datepicker.DatePicker
}

type Action struct {
	store TaskStateStore
}

func NewAction(store TaskStateStore) *Action {
	return &Action{store: store}
}

func (t *Action) CreateTaskFlow(ctx context.Context, params *TaskDTO) (*bot.SendMessageParams, error) {
	sID := strconv.FormatInt(params.ChatID, 10)
	if params.Init {
		state := TaskState{
			UserID:   params.UserID,
			NextStep: TitleStep,
			Data:     Task{},
		}
		err := t.store.Set(ctx, sID, &state)
		if err != nil {
			return nil, err
		}
		return &bot.SendMessageParams{
			ChatID: params.ChatID,
			Text:   "Enter the task title",
		}, nil
	}
	state, err := t.store.Get(ctx, sID)
	if err != nil {
		return nil, err
	}

	switch state.NextStep {
	case TitleStep:
		state.Data.Title = params.MsgText
		state.NextStep = DescStep
		err = t.store.Set(ctx, sID, state)
		if err != nil {
			return nil, err
		}
		return &bot.SendMessageParams{
			ChatID: params.ChatID,
			Text:   "Enter the task description",
		}, nil
	case DescStep:
		state.Data.Description = params.MsgText
		state.NextStep = DateTimeStep
		err = t.store.Set(ctx, sID, state)
		if err != nil {
			return nil, err
		}
		return &bot.SendMessageParams{
			ChatID:      params.ChatID,
			Text:        "Enter the task deadline date",
			ReplyMarkup: params.ReplyMarkup,
		}, nil
	case DateTimeStep:
		state.Data.DeadlineDate = params.Date
		state.NextStep = CompletedStepTaskStep
		err = t.store.Set(ctx, sID, state)
		return &bot.SendMessageParams{
			ChatID: params.ChatID,
			Text:   "Task was created",
		}, nil
	default:
		return nil, errors.New("not create task flow input")
	}
}
