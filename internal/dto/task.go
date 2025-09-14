package dto

import "time"

type CreateTaskInput struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	DeadlineDate time.Time `json:"deadline_date"`
}
