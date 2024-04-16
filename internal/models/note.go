package models

type Note struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Folder_id int    `json:"folder_id"`
}
