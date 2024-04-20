package models

type Task struct {
	Id           int    `json:"id"`
	Text         string `json:"text"`
	Description  string `json:"description"`
	Date         string `json:"date"`
	Is_completed bool   `json:"is_completed"`
	Task_id      *int   `json:"task_id"`
	Folder_id    int    `json:"folder_id"`
	Favourites   bool   `json:"favourites"`
	Folder_name  string `json:"folder_name"`
	Subtasks     []Task `json:"subtasks"`
}
