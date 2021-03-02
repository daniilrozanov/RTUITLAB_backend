package pkg

type ShopQuantity struct {
	ShopId int `json:"shop_id" db:"shop_id"`
	Quantity int `json:"quantity" db:"quantity"`
}