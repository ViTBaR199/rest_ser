package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/internal/models"

	_ "github.com/lib/pq"
)

type TaskRepositories interface {
	CreateTask(ctx context.Context, task models.Task) error
	DeleteTask(ctx context.Context, id_to_del int) error
	FetchTask(ctx context.Context, start, end int) ([]models.Task, error)
}

type taskRepositories struct {
	db *sql.DB
}

func NewTaskRepositories(db *sql.DB) TaskRepositories {
	return &taskRepositories{db: db}
}

func (r *taskRepositories) CreateTask(ctx context.Context, task models.Task) error {
	var taskID sql.NullInt64
	if task.Task_id != nil {
		taskID = sql.NullInt64{Int64: int64(*task.Task_id), Valid: true}
	} else {
		taskID = sql.NullInt64{Int64: 0, Valid: false}
	}
	_, err := r.db.ExecContext(ctx, "SELECT create_new_task($1, $2, $3, $4, $5)", task.Text, task.Is_completed, taskID, task.Folder_id, task.Favourites)
	return err
}

func (r *taskRepositories) DeleteTask(ctx context.Context, id_to_del int) error {
	_, err := r.db.ExecContext(ctx, "SELECT delete_line_task($1)", id_to_del)
	return err
}

func (r *taskRepositories) FetchTask(ctx context.Context, start, end int) ([]models.Task, error) {
	var result []models.Task

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM fetch_task($1, $2)", start, end)
	if err != nil {
		return nil, fmt.Errorf("querying fetch_task: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task
		var task_id sql.NullInt64

		if err := rows.Scan(&t.Id, &t.Text, &t.Is_completed, &task_id, &t.Folder_id, &t.Favourites); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}

		if task_id.Valid {
			temptask_id := int(task_id.Int64)
			t.Task_id = &temptask_id
		} else {
			t.Task_id = nil
		}

		result = append(result, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %v", err)
	}

	return result, nil
}
