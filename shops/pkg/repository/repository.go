package repository

type Authorization interface {
	ConfirmUser(name, password string) error
}

type Products interface {
	GetProductsByShop()
	GetProductsByCategory()
	GetProductsById()
	CreateProduct()
	UpdateProduct()
}

type Receipts interface {
	CreateReceipt()
	DeleteReceipt()
}

type Repository struct {
	Authorization
	Products
	Receipts
}
