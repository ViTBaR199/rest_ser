package models

type Task struct {
	Id           int    `json:"id"`
	Text         string `json:"text"`
	Is_completed bool   `json:"is_completed"`
	Task_id      int    `json:"task_id"`
	Folder_id    int    `json:"folder_id"`
}
