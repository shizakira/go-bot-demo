package usecase

import (
	"context"
	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/internal/dto"
)

type Repository interface {
	Add(ctx context.Context, t domain.Task) error
	GetAll(ctx context.Context) ([]*domain.Task, error)
}

type TaskService struct {
	repo Repository
}

func NewTaskService(repo Repository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, input dto.CreateTaskInput) error {
	task := domain.Task{
		UserID:       input.UserID,
		Title:        input.Title,
		Description:  input.Description,
		DeadlineDate: input.DeadlineDate,
	}
	return s.repo.Add(ctx, task)
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]*domain.Task, error) {
	return s.repo.GetAll(ctx)
}
