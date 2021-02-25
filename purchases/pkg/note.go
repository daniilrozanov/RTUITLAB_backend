package templates

import "time"

type Product struct {
	Id int `json:"id"`
	Title string `json:"title" binding:"required"`
	Date time.Time `json:"date"`
	Cost int `json:"cost"`
	Category string `json:"category"`
}

type UpdateProductInput struct {
	Title *string `json:"title"`
	Cost *int `json:"cost"`
	Category *string `json:"category"`
}