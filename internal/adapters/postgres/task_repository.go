package postgres

import (
	"context"
	"github.com/shizakira/daily-tg-bot/internal/domain"
	"time"
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
	_, err = stmt.ExecContext(ctx, t.UserID, t.Title, t.Description, t.Deadline)
	return err
}

func (tr *TaskRepository) getTasksByQuery(
	ctx context.Context,
	query string,
	params []any,
) ([]*domain.Task, error) {
	rows, err := tr.pool.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		newTask := new(domain.Task)
		if err = rows.Scan(
			&newTask.ID, &newTask.UserID, &newTask.Title, &newTask.Description, &newTask.Done,
			&newTask.Deadline, &newTask.CreatedAt, &newTask.ClosedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, newTask)
	}
	return tasks, err
}

func (tr *TaskRepository) CloseTask(ctx context.Context, id int64, isDone bool) error {
	stmt, err := tr.pool.PrepareContext(ctx,
		"update tasks set closed_at = $1, done = $2 where id = $3",
	)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, time.Now().Format("2006-01-02 15:04:05"), isDone, id)
	return err
}

func (tr *TaskRepository) GetOpenByUserID(ctx context.Context, userId int64) ([]*domain.Task, error) {
	return tr.getTasksByQuery(
		ctx,
		"select * from tasks where user_id = $1 and done = false and closed_at is null",
		[]any{userId},
	)
}

func (tr *TaskRepository) GetExpiredTasks(ctx context.Context) ([]*domain.Task, error) {
	expired, err := tr.getTasksByQuery(ctx, `
		SELECT *
		FROM tasks 
		WHERE deadline < NOW()
			AND done = false
			AND closed_at IS NULL;
    `, []any{})
	if err != nil {
		return nil, err
	}
	return expired, nil
}

func (tr *TaskRepository) GetSoonExpiredTasks(ctx context.Context) ([]*domain.Task, error) {
	soon, err := tr.getTasksByQuery(ctx, `
        SELECT * 
		FROM tasks
		WHERE deadline BETWEEN NOW() AND NOW() + INTERVAL '1 hour'
          AND done = false
          AND closed_at IS NULL
    `, []any{})
	if err != nil {
		return nil, err
	}

	return soon, nil
}
