package pkg

type ShopQuantityJSON struct {
	ShopId		int `json:"shop_id" db:"shop_id"`
	Quantity	int `json:"quantity" db:"quantity"`
}

type ProductJSON struct {
	ID          int					`json:"id"`
	Title       string				`json:"title"`
	Cost        int					`json:"cost"`
	Quantity	int					`json:"quantity"`
	Category    string				`json:"category"`
}

type CartJSON struct {
	Shop		Shop 			`json:"shop"`
	Products	[]ProductJSON	`json:"products"`
	SummaryCost	int				`json:"summary_cost"`
	Index		int				`json:"index"`
}

type CartItemsOnDeleteJSON struct {
	 ShopID		int `json:"shop_id" binding:"required"`
	 ProductID	int `json:"product_id" binding:"required"`
	 Quantity	int `json:"quantity"`
}