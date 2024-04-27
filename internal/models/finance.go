package models

type Finance struct {
	Id            int    `json:"id"`
	CategoryName  string `json:"category_name"`
	CategoryPhoto string `json:"category_photo"`
	CategoryColor int    `json:"category_color"`
	YearMonth     string `json:"yearMonth"`
	Price         int    `json:"price"`
	Date          string `json:"date"`
	Currency      string `json:"currency"`
	Folder_id     int    `json:"folder_id"`
}
