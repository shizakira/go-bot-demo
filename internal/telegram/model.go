package telegram

import "time"

type TaskStep string

const (
	TitleStep             TaskStep = "title"
	DescStep              TaskStep = "description"
	DateTimeStep          TaskStep = "datetime"
	CompletedStepTaskStep TaskStep = "completed"
)

type TaskState struct {
	UserID   int64    `json:"user_id"`
	NextStep TaskStep `json:"next_step"`
	Data     Task     `json:"data"`
}

type Task struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	DeadlineDate time.Time `json:"deadline_date"`
}
