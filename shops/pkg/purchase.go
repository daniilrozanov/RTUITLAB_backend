package pkg

import "time"

type Purchase struct {
	Id int `json:"id"`
	Title string `json:"title" binding:"required"`
	Date time.Time `json:"date"`
	Cost int `json:"cost"`
	Category string `json:"category"`
}
