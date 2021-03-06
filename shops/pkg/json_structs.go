package pkg

import "time"

type ShopQuantityJSON struct {
	ShopId   int `json:"shop_id" db:"shop_id"`
	Quantity int `json:"quantity" db:"quantity"`
}

type ProductJSON struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Cost       int    `json:"cost"`
	Quantity   int    `json:"quantity"`
	EntireCost int    `json:"entire_cost"`
	Category   string `json:"category"`
}

type CartJSON struct {
	Shop        Shop          `json:"shop"`
	Products    []ProductJSON `json:"products"`
	SummaryCost int           `json:"summary_cost"`
}

type CartItemsOnDeleteJSON struct {
	ShopID    int `json:"shop_id" binding:"required"`
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity"`
}

type ReceiptJSON struct {
	CartJSON
	PayOption   string `json:"pay_option"`
	CreatedTime time.Time
}

type UserReceiptMapJSON struct {
	Receipt ReceiptJSON `json:"receipt"`
	UserID  int         `json:"user_id"`
}

type CartItemJSON struct {
	CartItem
	Category string `json:",omitempty"`
}

type CreateProductData struct {
	Prod       Product         `json:"product" binding:"required"`
	ShopsCount []ShopsProducts `json:"map" binding:"required"`
}