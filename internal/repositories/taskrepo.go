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
	FetchTask(ctx context.Context, start, end int, folder_id ...int) ([]models.Task, error)
	UpdateTask(ctx context.Context, task models.Task) error
	CountTask(ctx context.Context) (int, error)
	CountTaskFavourites(ctx context.Context) (int, error)
	FetchTaskFavourites(ctx context.Context, start, end int, folder_id ...int) ([]models.Task, error)
}

type taskRepositories struct {
	db *sql.DB
}

func NewTaskRepositories(db *sql.DB) TaskRepositories {
	return &taskRepositories{db: db}
}

func (r *taskRepositories) CreateTask(ctx context.Context, task models.Task) error {
	var tasks []models.Task
	var err error
	tasks, err = r.FetchTask(ctx, 0, 0)

	if err != nil {
		return fmt.Errorf("an empty task list")
	}

	if len(tasks) == 0 {
		return fmt.Errorf("no tasks available for checking")
	}

	for _, t := range tasks {
		if task.Task_id != nil && task.Id == *t.Task_id {
			err := fmt.Errorf("the parent element cannot be a child")
			return err
		} else if t.Task_id != nil && *task.Task_id == t.Id {
			err := fmt.Errorf("an element cannot be a child of another child element")
			return err
		}
	}

	var taskID sql.NullInt64
	if task.Task_id != nil {
		taskID = sql.NullInt64{Int64: int64(*task.Task_id), Valid: true}
	} else {
		taskID = sql.NullInt64{Int64: 0, Valid: false}
	}
	_, err = r.db.ExecContext(ctx, "SELECT create_new_task($1, $2, $3, $4, $5, $6, $7)",
		task.Text, task.Description, task.Is_completed, taskID, task.Folder_id, task.Favourites, task.Date)
	return err
}

func (r *taskRepositories) DeleteTask(ctx context.Context, id_to_del int) error {
	_, err := r.db.ExecContext(ctx, "SELECT delete_line_task($1)", id_to_del)
	return err
}

func (r *taskRepositories) FetchTask(ctx context.Context, start, end int, folder_id ...int) ([]models.Task, error) {
	taskMap := make(map[int]*models.Task)
	allTasks := []models.Task{} // Список для задач без task_id

	var rows *sql.Rows
	var err error

	// Выбор запроса в зависимости от наличия folder_id
	if start == 0 && end == 0 {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_task()")
	} else if len(folder_id) > 0 {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_task($1, $2, $3)", start, end, folder_id[0])
	} else {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_task($1, $2)", start, end)
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
	//проверить, есть ли у него родительский task
	////тогда родителя можно заменить на другого (кроме id, который сам является дочерним)
	//проверить, является ли он чьим-то родителем
	////тогда он не может иметь родителя кроме null

	tasks, err := r.FetchTask(ctx, 0, 0)

	if err != nil {
		return fmt.Errorf("an empty task list")
	}

	if len(tasks) == 0 {
		return fmt.Errorf("no tasks available for checking")
	}

	for _, t := range tasks {
		if task.Task_id != nil && task.Id == *t.Task_id {
			err := fmt.Errorf("the parent element cannot be a child")
			return err
		} else if t.Task_id != nil && *task.Task_id == t.Id {
			err := fmt.Errorf("an element cannot be a child of another child element")
			return err
		}
	}

	_, err = r.db.ExecContext(ctx, "SELECT update_task($1, $2, $3, $4, $5, $6, $7, $8)", task.Id, task.Text, task.Description, task.Date, task.Is_completed,
		task.Favourites, task.Task_id, task.Folder_id)

	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}
	return err
}

func (r *taskRepositories) CountTask(ctx context.Context) (int, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM count_task()")
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

func (r *taskRepositories) CountTaskFavourites(ctx context.Context) (int, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM count_task_favourites()")
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

func (r *taskRepositories) FetchTaskFavourites(ctx context.Context, start, end int, folder_id ...int) ([]models.Task, error) {
	taskMap := make(map[int]*models.Task)
	allTasks := []models.Task{} // Список для задач без task_id

	var rows *sql.Rows
	var err error

	// Выбор запроса в зависимости от наличия folder_id
	if len(folder_id) > 0 {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_task_favourites($1, $2, $3)", start, end, folder_id[0])
	} else {
		rows, err = r.db.QueryContext(ctx, "SELECT * FROM fetch_task_favourites($1, $2)", start, end)
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
