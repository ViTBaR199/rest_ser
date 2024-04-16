package models

import "time"

type Finance struct {
	Id        int       `json:"id"`
	Price     int       `json:"price"`
	Date      time.Time `json:"date"`
	Currency  string    `json:"currency"`
	Folder_id int       `json:"folder_id"`
}
