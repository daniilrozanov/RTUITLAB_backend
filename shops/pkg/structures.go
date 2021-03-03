package pkg

import (
	"time"
)

type Receip struct {
	Id				int		`json:"id"`
	ShopID			int		`json:"shop_id"`
	UserID			int		`json:"user_id"`
	CartItemsNumber	int		`json:""`
	PayCase			string	`json:"pay_case"`
	CreateDate		time.Time
}

type Shop struct {
	ID		int   	`json:"id"`
	Title	string	`json:"title"`
	Address	string	`json:"address"`
	Phone	string	`json:"number"`
}

type Product struct {
	ID          int                `json:"id"`
	Title       string             `json:"title" binding:"required"`
	Description string             `json:"description" binding:"required"`
	Cost        int                `json:"cost" binding:"required"`
	Category    string             `json:"category" binding:"required"`
	Code        int                `json:"code" binding:"required"`
	SQ          []ShopQuantityJSON `json:"map"`
}

type ShopsProducts struct {
	ID        int `json:"id"`
	ShopID    int `json:"shop_id" binding:"required"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity" binding:"required"`
}

type CartItem struct {
	ID        	int `json:"id"`
	ProductID 	int `json:"product_id" binding:"required" db:"product_id"`
	ShopID    	int `json:"shop_id" binding:"required" db:"shop_id"`
	UserID    	int `json:"user_id" db:"user_id"`
	Quantity	int `json:"quantity" db:"quantity"`
	CartNumber	int `json:"cart_number" db:"index"`
}

type UserCarts struct {
	ID				int `json:"id"`
	UserID			int
	ShopID			int `db:"shop_id"`
	NumberOfCarts	int `db:"carts_number"`
}