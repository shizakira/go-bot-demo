package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/internal/dto"
)

type taskRepositoryMock struct {
	addedTasks  []domain.Task
	addErr      error
	returned    []*domain.Task
	getErr      error
	closedTasks []closeCall
	closeErr    error
}

type closeCall struct {
	id   int64
	done bool
}

func (m *taskRepositoryMock) Add(_ context.Context, task domain.Task) error {
	m.addedTasks = append(m.addedTasks, task)
	return m.addErr
}

func (m *taskRepositoryMock) GetOpenByUserID(_ context.Context, _ int64) ([]*domain.Task, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.returned, nil
}

func (m *taskRepositoryMock) CloseTask(_ context.Context, id int64, isDone bool) error {
	m.closedTasks = append(m.closedTasks, closeCall{id: id, done: isDone})
	return m.closeErr
}

func (m *taskRepositoryMock) GetExpiredTasks(context.Context) ([]*domain.Task, error) {
	return nil, nil
}

func (m *taskRepositoryMock) GetSoonExpiredTasks(context.Context) ([]*domain.Task, error) {
	return nil, nil
}

func TestTaskService_CreateTask(t *testing.T) {
	repo := &taskRepositoryMock{}
	service := NewTaskService(repo)
	deadline := time.Now().Add(2 * time.Hour)
	input := dto.CreateTaskInput{
		UserID:       42,
		Title:        "Test title",
		Description:  "Test description",
		DeadlineDate: deadline,
	}

	if err := service.CreateTask(context.Background(), input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.addedTasks) != 1 {
		t.Fatalf("expected 1 task added, got %d", len(repo.addedTasks))
	}

	added := repo.addedTasks[0]
	if added.UserID != input.UserID {
		t.Errorf("expected user id %d, got %d", input.UserID, added.UserID)
	}
	if added.Title != input.Title {
		t.Errorf("expected title %q, got %q", input.Title, added.Title)
	}
	if added.Description != input.Description {
		t.Errorf("expected description %q, got %q", input.Description, added.Description)
	}
	if !added.Deadline.Equal(deadline) {
		t.Errorf("expected deadline %v, got %v", deadline, added.Deadline)
	}
}

func TestTaskService_GetOpenTasksByUserID(t *testing.T) {
	returnedTasks := []*domain.Task{{ID: 1}, {ID: 2}}
	repo := &taskRepositoryMock{returned: returnedTasks}
	service := NewTaskService(repo)

	output, err := service.GetOpenTasksByUserID(context.Background(), dto.GetAllTasksByUserIdInput{UserID: 99})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(output.Tasks) != len(returnedTasks) {
		t.Fatalf("expected %d tasks, got %d", len(returnedTasks), len(output.Tasks))
	}
	for i, task := range output.Tasks {
		if task.ID != returnedTasks[i].ID {
			t.Errorf("expected task ID %d, got %d", returnedTasks[i].ID, task.ID)
		}
	}
}

func TestTaskService_GetOpenTasksByUserID_Error(t *testing.T) {
	repo := &taskRepositoryMock{getErr: errors.New("db failure")}
	service := NewTaskService(repo)

	output, err := service.GetOpenTasksByUserID(context.Background(), dto.GetAllTasksByUserIdInput{UserID: 7})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if output.Tasks != nil {
		t.Fatalf("expected no tasks on error, got %v", output.Tasks)
	}
}

func TestTaskService_CloseTask(t *testing.T) {
	repo := &taskRepositoryMock{}
	service := NewTaskService(repo)

	input := dto.CloseTaskInput{TaskID: 123, IsDone: true}
	if err := service.CloseTask(context.Background(), input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.closedTasks) != 1 {
		t.Fatalf("expected 1 close call, got %d", len(repo.closedTasks))
	}
	call := repo.closedTasks[0]
	if call.id != input.TaskID {
		t.Errorf("expected task id %d, got %d", input.TaskID, call.id)
	}
	if call.done != input.IsDone {
		t.Errorf("expected done %v, got %v", input.IsDone, call.done)
	}
}
