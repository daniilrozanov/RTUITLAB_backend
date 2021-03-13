package pkg

type FabricProductJSON struct {
	CreateProductData
	Power int
}

type CreateProductData struct {
	Prod       Product         `json:"product" binding:"required"`
	ShopsCount []ShopsProducts `json:"map" binding:"required"`
}

type Product struct {
	Title       string             `json:"title" binding:"required"`
	Description string             `json:"description" binding:"required"`
	Cost        int                `json:"cost" binding:"required"`
	Category    string             `json:"category" binding:"required"`
	Code        int                `json:"code" binding:"required"`
}

type ShopsProducts struct {
	ShopID    int `json:"shop_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}
/*
type ProductJSON struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Cost       int    `json:"cost"`
	Quantity   int    `json:"quantity"`
	EntireCost int    `json:"entire_cost"`
	Category   string `json:"category"`
}

type Shop struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Address string `json:"address"`
	Phone   string `json:"number"`
}

type CreateProductData struct {
	Prod       pkg.Product         `json:"product" binding:"required"`
	ShopsCount []pkg.ShopsProducts `json:"map" binding:"required"`
}
*/