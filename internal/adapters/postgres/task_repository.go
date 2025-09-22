package postgres

import (
	"context"
	"github.com/shizakira/daily-tg-bot/internal/domain"
)

type TaskRepository struct {
	pool *Pool
}

func NewTaskRepository(pool *Pool) *TaskRepository {
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

func (tr *TaskRepository) GetAllByUserID(ctx context.Context, userId int64) ([]*domain.Task, error) {
	rows, err := tr.pool.QueryContext(ctx, "select * from tasks where user_id = $1", userId)
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
