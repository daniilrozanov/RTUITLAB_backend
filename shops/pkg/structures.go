package pkg

import (
	"time"
)

type Receip struct {
	Id uint `json:"id" gorm:"primaryKey;autoIncrement"`
	ShopID uint
	UserID uint
	CartItemsNumber uint
	CreateDate time.Time
}

type Shop struct {
	ID      int   	`json:"id"`
	Title   string	`json:"title"`
	Address string	`json:"address"`
	Phone  string	`json:"number"`
}

type Product struct {
	ID          int            `json:"id"`
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Cost        int            `json:"cost" binding:"required"`
	Category    string         `json:"category" binding:"required"`
	Code        int            `json:"code" binding:"required"`
	SQ          []ShopQuantity `json:"map"`
}

type ShopsProducts struct {
	ID        int `json:"id"`
	ShopID    int `json:"shop_id" binding:"required"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity" binding:"required"`
}

type CartItem struct {
	ID int `json:"id"`
	ProductID int ``
	ShopID int ``
	UserID int ``
	Count int ``
	Number int
}

type UsersCarts struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID int
	ShopID int
	NumberOfCarts int
}
