package usecase

import (
	"context"
	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/internal/dto"
)

type TaskRepository interface {
	Add(ctx context.Context, t domain.Task) error
	GetOpenByUserID(ctx context.Context, userUd int64) ([]*domain.Task, error)
	CloseTask(ctx context.Context, id int64, isDone bool) error
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, input dto.CreateTaskInput) error {
	task := domain.Task{
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Description,
		Deadline:    input.DeadlineDate,
	}
	return s.repo.Add(ctx, task)
}

func (s *TaskService) GetOpenTasksByUserID(
	ctx context.Context,
	input dto.GetAllTasksByUserIdInput,
) (dto.GetAllTasksByUserIdOutput, error) {
	output := dto.GetAllTasksByUserIdOutput{}
	tasks, err := s.repo.GetOpenByUserID(ctx, input.UserID)
	if err == nil {
		output.Tasks = tasks
	}
	return output, err
}

func (s *TaskService) CloseTask(ctx context.Context, input dto.CloseTaskInput) error {
	return s.repo.CloseTask(ctx, input.TaskID, input.IsDone)
}

func (s *TaskService) SendNotifyForExpiredTasks() {

}

func (s *TaskService) SendNotifyFortNearExpiredTasks() {

}
