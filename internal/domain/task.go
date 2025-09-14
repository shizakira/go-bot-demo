package domain

import "time"

type Task struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	DeadlineDate time.Time `json:"deadline_date"`
}
