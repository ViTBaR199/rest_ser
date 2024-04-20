package models

type Finance struct {
	Id        int    `json:"id"`
	Price     int    `json:"price"`
	Date      string `json:"date"`
	Currency  string `json:"currency"`
	Folder_id int    `json:"folder_id"`
}
