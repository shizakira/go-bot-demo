package postgres

import (
	"context"
	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/pkg/database"
)

type TaskRepository struct {
	pool *database.PostgresPool
}

func NewTaskRepository(pool *database.PostgresPool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (tr *TaskRepository) Add(ctx context.Context, t domain.Task) error {
	stmt, err := tr.pool.PrepareContext(ctx,
		"insert into tasks (user_id, title, description, deadline) values ($1, $2, $3, $4)",
	)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, t.UserID, t.Title, t.Description, t.DeadlineDate)
	return err
}

func (tr *TaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	rows, err := tr.pool.QueryContext(ctx, "select * from tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task

	for rows.Next() {
		newTask := new(domain.Task)
		if err = rows.Scan(&newTask.ID, &newTask.UserID, &newTask.Title, &newTask.Description, &newTask.DeadlineDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, newTask)
	}
	return tasks, err
}
