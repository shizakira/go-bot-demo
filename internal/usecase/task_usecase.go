package usecase

import (
	"context"
	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/internal/dto"
)

type TaskRepository interface {
	Add(ctx context.Context, t domain.Task) error
	GetAllByUserID(ctx context.Context, userUd int64) ([]*domain.Task, error)
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
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

func (s *TaskService) GetAllTasksByUserID(
	ctx context.Context,
	input dto.GetAllTasksByUserIdInput,
) (dto.GetAllTasksByUserIdOutput, error) {
	output := dto.GetAllTasksByUserIdOutput{}
	tasks, err := s.repo.GetAllByUserID(ctx, input.UserID)
	if err == nil {
		output.Tasks = tasks
	}
	return output, err
}
