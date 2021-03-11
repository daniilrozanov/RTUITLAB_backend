package buisness

import (
	"purchases/pkg"
	"purchases/pkg/repository"
)

type Authorization interface {
	CreateUser(user templates.User) (int, error)
	GenerateToken(username, password string) (string, error)
	GetUserId(username, password string) (int, error)
	ParseToken(token string) (int, error)
}

type ProductLogging interface {
	CreateProduct(userId int, prod *templates.Product) (int, error)
	GetAllProducts(userId int) ([]templates.Product, error)
	GetProductById(userId, prodId int) (templates.Product, error)
	UpdateProduct(userId, prodId int, input *templates.UpdateProductInput) error
	DeleteProduct(userId, prodId int) error
}

type Synchronization interface {
	StartConsume() error
	GetReceipts (userID int) ([]templates.ReceiptJSON, error)
}

type Buisness struct {
	Authorization
	ProductLogging
	Synchronization
}

func InitBuisnessLayer(r *repository.Repository, rabbit *RabbitStruct) *Buisness {
	return &Buisness{
		Authorization: NewAuthService(*r),
		ProductLogging: NewProductService(*r),
		Synchronization: NewSynchroService(*rabbit, *r),
	}
}