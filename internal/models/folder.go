package models

type Folder struct {
	ID           *int   `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Image        string `json:"image"`
	Color        *int   `json:"color"`
	Folder_count int    `json:"folder_count"`
	User_id      *int   `json:"user_id"`
}
