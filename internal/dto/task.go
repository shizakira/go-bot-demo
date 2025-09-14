package dto

import "time"

type CreateTaskInput struct {
	UserID       int64     `json:"user_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	DeadlineDate time.Time `json:"deadline_date"`
}
