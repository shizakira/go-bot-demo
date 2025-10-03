package telegram

import (
	"context"
	"testing"

	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/internal/usecase"
)

type handleTaskRepoMock struct {
	closedCalls []closeArgs
	closeErr    error
}

type closeArgs struct {
	id   int64
	done bool
}

func (m *handleTaskRepoMock) Add(context.Context, domain.Task) error {
	return nil
}

func (m *handleTaskRepoMock) GetOpenByUserID(context.Context, int64) ([]*domain.Task, error) {
	return nil, nil
}

func (m *handleTaskRepoMock) CloseTask(_ context.Context, id int64, isDone bool) error {
	m.closedCalls = append(m.closedCalls, closeArgs{id: id, done: isDone})
	return m.closeErr
}

func (m *handleTaskRepoMock) GetExpiredTasks(context.Context) ([]*domain.Task, error) {
	return nil, nil
}

func (m *handleTaskRepoMock) GetSoonExpiredTasks(context.Context) ([]*domain.Task, error) {
	return nil, nil
}

func TestGetValueFromQueryByRe(t *testing.T) {
	tb := &Bot{}
	value, err := tb.getValueFromQueryByRe(`ID: ([0-9]+)`, "ID: 15\nTitle: task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "15" {
		t.Fatalf("expected value 15, got %s", value)
	}
}

func TestGetValueFromQueryByRe_NotFound(t *testing.T) {
	tb := &Bot{}
	if _, err := tb.getValueFromQueryByRe(`ID: ([0-9]+)`, "no identifier"); err == nil {
		t.Fatal("expected error when regex does not match")
	}
}

func TestHandleTaskClosure(t *testing.T) {
	repo := &handleTaskRepoMock{}
	service := usecase.NewTaskService(repo)
	b := &Bot{taskService: service}

	query := "ID: 21\nTitle: test"
	if err := b.handleTaskClosure(context.Background(), query, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.closedCalls) != 1 {
		t.Fatalf("expected 1 close call, got %d", len(repo.closedCalls))
	}
	call := repo.closedCalls[0]
	if call.id != 21 {
		t.Fatalf("expected task id 21, got %d", call.id)
	}
	if !call.done {
		t.Fatalf("expected done flag true")
	}
}

func TestHandleTaskClosure_InvalidQuery(t *testing.T) {
	repo := &handleTaskRepoMock{}
	service := usecase.NewTaskService(repo)
	b := &Bot{taskService: service}

	if err := b.handleTaskClosure(context.Background(), "Title only", false); err == nil {
		t.Fatal("expected error for query without id")
	}
	if len(repo.closedCalls) != 0 {
		t.Fatalf("expected no close calls, got %d", len(repo.closedCalls))
	}
}

func TestHandleTaskClosure_InvalidID(t *testing.T) {
	repo := &handleTaskRepoMock{}
	service := usecase.NewTaskService(repo)
	b := &Bot{taskService: service}

	if err := b.handleTaskClosure(context.Background(), "ID: abc", false); err == nil {
		t.Fatal("expected error for invalid id")
	}
	if len(repo.closedCalls) != 0 {
		t.Fatalf("expected no close calls, got %d", len(repo.closedCalls))
	}
}
