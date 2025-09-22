package telegram

import "github.com/shizakira/daily-tg-bot/internal/dto"

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
