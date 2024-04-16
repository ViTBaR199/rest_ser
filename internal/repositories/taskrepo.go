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
	FetchTask(ctx context.Context, start, end int) ([][]string, error)
}

type taskRepositories struct {
	db *sql.DB
}

func NewTaskRepositories(db *sql.DB) TaskRepositories {
	return &taskRepositories{db: db}
}

func (r *taskRepositories) CreateTask(ctx context.Context, task models.Task) error {
	var taskID sql.NullInt64
	if task.Task_id != 0 {
		taskID = sql.NullInt64{Int64: int64(task.Task_id), Valid: true}
	}
	_, err := r.db.ExecContext(ctx, "SELECT create_new_task($1, $2, $3, $4)", task.Text, task.Is_completed, taskID, task.Folder_id)
	return err
}

func (r *taskRepositories) DeleteTask(ctx context.Context, id_to_del int) error {
	_, err := r.db.ExecContext(ctx, "SELECT delete_line_task($1)", id_to_del)
	return err
}

func (r *taskRepositories) FetchTask(ctx context.Context, start, end int) ([][]string, error) {
	var result [][]string

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM fetch_task($1, $2)", start, end)
	if err != nil {
		return nil, fmt.Errorf("querying fetch_task: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, folder_id int
		var task_id sql.NullInt32
		var text string
		var is_completed bool
		if err := rows.Scan(&id, &text, &is_completed, &task_id, &folder_id); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}

		taskIdStr := "NULL"
		if task_id.Valid {
			taskIdStr = fmt.Sprintf("%d", task_id.Int32)
		}

		rowData := []string{fmt.Sprintf("%d", id), text, fmt.Sprintf("%t", is_completed), taskIdStr, fmt.Sprintf("%d", folder_id)}
		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %v", err)
	}

	return result, nil
}
