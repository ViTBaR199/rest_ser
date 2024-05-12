package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/internal/models"

	_ "github.com/lib/pq"
)

type TaskRepositories interface {
	CreateTask(ctx context.Context, task models.Task) (int, error)
	DeleteTask(ctx context.Context, id_to_del int) error
	FetchTask(ctx context.Context, user_id, start, end int, folder_id ...int) ([]models.Task, error)
	UpdateTask(ctx context.Context, task models.Task) error
	CountTask(ctx context.Context, user_id int) (int, error)
	CountTaskFavourites(ctx context.Context, user_id int) (int, error)
	FetchTaskFavourites(ctx context.Context, user_id, start, end int, folder_id ...int) ([]models.Task, error)
	GetUserByTask(taskID int) (int, error)
}

type taskRepositories struct {
	db *sql.DB
}

func NewTaskRepositories(db *sql.DB) TaskRepositories {
	return &taskRepositories{db: db}
}

func (r *taskRepositories) CreateTask(ctx context.Context, task models.Task) (int, error) {
	var err error
	var taskID sql.NullInt64
	if task.Task_id != nil {
		taskID = sql.NullInt64{Int64: int64(*task.Task_id), Valid: true}
		var parentIsChild bool
		err = r.db.QueryRowContext(ctx, "SELECT checking_for_childishness($1)", task.Task_id).Scan(&parentIsChild)
		if err != nil {
			return 0, err
		}
		if parentIsChild {
			return 0, fmt.Errorf("cannot create a subtask for a task that is already a subtask")
		}
	} else {
		taskID = sql.NullInt64{Int64: 0, Valid: false}
	}
	var newId int
	err = r.db.QueryRowContext(ctx, "SELECT create_new_task($1, $2, $3, $4, $5, $6, $7)",
		task.Text, task.Description, task.Is_completed, taskID, task.Folder_id, task.Favourites, task.Date).Scan(&newId)
	return newId, nil
}

func (r *taskRepositories) DeleteTask(ctx context.Context, id_to_del int) error {
	_, err := r.db.ExecContext(ctx, "SELECT delete_line_task($1)", id_to_del)
	return err
}

func (r *taskRepositories) FetchTask(ctx context.Context, user_id, start, end int, folder_id ...int) ([]models.Task, error) {
	taskMap := make(map[int]*models.Task)
	allTasks := []models.Task{} // Список для задач без task_id

	var rows *sql.Rows
	var err error

	// Выбор запроса в зависимости от наличия folder_id
	if start == 0 && end == 0 {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_task($1)", user_id)
	} else if len(folder_id) > 0 {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_task($1, $2, $3, $4)", user_id, start, end, folder_id[0])
	} else {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_task($1, $2, $3)", user_id, start, end)
	}

	if err != nil {
		return nil, fmt.Errorf("querying fetch_task: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task
		var task_id sql.NullInt64
		var description sql.NullString

		if err := rows.Scan(&t.Id, &t.Text, &description, &t.Date, &t.Is_completed, &t.Favourites, &task_id, &t.Folder_id, &t.Folder_name); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}

		if task_id.Valid {
			temptask_id := int(task_id.Int64)
			t.Task_id = &temptask_id
		}

		t.Description = description.String

		taskMap[t.Id] = &t
		allTasks = append(allTasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %v", err)
	}

	// После создания всех задач, организуем подзадачи
	for _, task := range allTasks {
		if task.Task_id != nil && *task.Task_id != task.Id {
			if parent, exists := taskMap[*task.Task_id]; exists {
				parent.Subtasks = append(parent.Subtasks, task)
			}
		}
	}

	// Формирование конечного списка результатов, добавляем только родительские задачи
	var result []models.Task
	for _, task := range taskMap {
		if task.Task_id == nil { // Добавляем только корневые задачи
			result = append(result, *task)
		}
	}

	return result, nil
}

func (r *taskRepositories) UpdateTask(ctx context.Context, task models.Task) error {
	var err error
	var canUpdate bool

	// Проверяем возможность обновления
	err = r.db.QueryRowContext(ctx, "SELECT check_task_constraints($1, $2)", task.Id, task.Task_id).Scan(&canUpdate)
	if err != nil {
		return fmt.Errorf("error checking task constraints: %v", err)
	}

	if !canUpdate {
		return fmt.Errorf("task constraints prevent updating")
	}

	_, err = r.db.ExecContext(ctx, "SELECT update_task($1, $2, $3, $4, $5, $6, $7, $8)", task.Id, task.Text, task.Description, task.Date, task.Is_completed,
		task.Favourites, task.Task_id, task.Folder_id)

	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}
	return err
}

func (r *taskRepositories) CountTask(ctx context.Context, user_id int) (int, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM count_task($1)", user_id)
	if err != nil {
		return 0, fmt.Errorf("querying count_task: %v", err)
	}
	defer rows.Close()

	var count sql.NullInt64
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, fmt.Errorf("scanning row: %v", err)
		}
		if !count.Valid {
			// Обрабатываем случай, когда получено NULL значение
			return 0, fmt.Errorf("received NULL value for count")
		}
		return int(count.Int64), nil
	} else {
		// Обрабатываем случай, когда не возвращено ни одной строки
		if err := rows.Err(); err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("no rows returned")
	}
}

func (r *taskRepositories) CountTaskFavourites(ctx context.Context, user_id int) (int, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM count_task_favourites($1)", user_id)
	if err != nil {
		return 0, fmt.Errorf("querying count_task_favourites: %v", err)
	}
	defer rows.Close()

	var count sql.NullInt64
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, fmt.Errorf("scanning row: %v", err)
		}
		if !count.Valid {
			// Обрабатываем случай, когда получено NULL значение
			return 0, fmt.Errorf("received NULL value for count")
		}
		return int(count.Int64), nil
	} else {
		// Обрабатываем случай, когда не возвращено ни одной строки
		if err := rows.Err(); err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("no rows returned")
	}
}

func (r *taskRepositories) FetchTaskFavourites(ctx context.Context, user_id, start, end int, folder_id ...int) ([]models.Task, error) {
	taskMap := make(map[int]*models.Task)
	allTasks := []models.Task{} // Список для задач без task_id

	var rows *sql.Rows
	var err error

	// Выбор запроса в зависимости от наличия folder_id
	if len(folder_id) > 0 {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_task_favourites($1, $2, $3, $4)", user_id, start, end, folder_id[0])
	} else {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_task_favourites($1, $2, $3)", user_id, start, end)
	}

	if err != nil {
		return nil, fmt.Errorf("querying fetch_task: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task
		var task_id sql.NullInt64
		var description sql.NullString

		if err := rows.Scan(&t.Id, &t.Text, &description, &t.Date, &t.Is_completed, &t.Favourites, &task_id, &t.Folder_id, &t.Folder_name); err != nil {
			return nil, fmt.Errorf("scanning row: %v", err)
		}

		if task_id.Valid {
			temptask_id := int(task_id.Int64)
			t.Task_id = &temptask_id
		}

		t.Description = description.String

		taskMap[t.Id] = &t
		allTasks = append(allTasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %v", err)
	}

	// После создания всех задач, организуем подзадачи
	for _, task := range allTasks {
		if task.Task_id != nil && *task.Task_id != task.Id {
			if parent, exists := taskMap[*task.Task_id]; exists {
				parent.Subtasks = append(parent.Subtasks, task)
			}
		}
	}

	// Формирование конечного списка результатов, добавляем только родительские задачи
	var result []models.Task
	for _, task := range taskMap {
		if task.Task_id == nil { // Добавляем только корневые задачи
			result = append(result, *task)
		}
	}

	return result, nil
}

func (r *taskRepositories) GetUserByTask(taskID int) (int, error) {
	var userId int
	err := r.db.QueryRow("SELECT get_user_by_task($1)", taskID).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
