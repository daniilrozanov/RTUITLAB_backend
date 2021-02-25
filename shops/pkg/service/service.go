package service

type Authorization interface {
	ConfirmUser(name, password string) (int, error)
	GenerateToken(id int) (string, error)
	ParseToken(token string) (int, error)
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

type Service struct {
	Authorization
	Products
	Receipts
}

func InitNewService(authorization Authorization) *Service {
	return &Service{Authorization: authorization}
}


