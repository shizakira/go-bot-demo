package dto

import (
	"github.com/shizakira/daily-tg-bot/internal/domain"
	"time"
)

type CreateTaskInput struct {
	UserID       int64     `json:"user_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	DeadlineDate time.Time `json:"deadline_date"`
}

type GetAllTasksByUserIdInput struct {
	UserID int64
}

type GetAllTasksByUserIdOutput struct {
	Tasks []*domain.Task
}

type CloseTaskInput struct {
	TaskID int64
	IsDone bool
}
