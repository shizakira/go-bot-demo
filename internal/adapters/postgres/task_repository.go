package postgres

import (
	"context"
	"time"

	"github.com/shizakira/daily-tg-bot/internal/domain"
)

type TaskRepository struct {
	pool *Pool
}

func NewTaskRepository(pool *Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (tr *TaskRepository) Add(ctx context.Context, t domain.Task) error {
	_, err := tr.pool.ExecContext(
		ctx,
		"insert into tasks (user_id, title, description, deadline) values ($1, $2, $3, $4)",
		t.UserID,
		t.Title,
		t.Description,
		t.Deadline,
	)
	return err
}

func (tr *TaskRepository) getTasksByQuery(
	ctx context.Context,
	query string,
	params ...any,
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (tr *TaskRepository) CloseTask(ctx context.Context, id int64, isDone bool) error {
	_, err := tr.pool.ExecContext(
		ctx,
		"update tasks set closed_at = $1, done = $2 where id = $3",
		time.Now().UTC(),
		isDone,
		id,
	)
	return err
}

func (tr *TaskRepository) GetOpenByUserID(ctx context.Context, userId int64) ([]*domain.Task, error) {
	return tr.getTasksByQuery(
		ctx,
		"select * from tasks where user_id = $1 and done = false and closed_at is null",
		userId,
	)
}

func (tr *TaskRepository) GetExpiredTasks(ctx context.Context) ([]*domain.Task, error) {
	return tr.getTasksByQuery(ctx, `
		select *
		from tasks 
		where deadline < now()
			and done = false
			and closed_at is null;
    `)
}

func (tr *TaskRepository) GetSoonExpiredTasks(ctx context.Context) ([]*domain.Task, error) {
	return tr.getTasksByQuery(ctx, `
        select * 
		from tasks
		where deadline between now() and now() + interval '1 hour'
          and done = false
          and closed_at is null
    `)
}
