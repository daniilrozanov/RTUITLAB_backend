package service

import (
	"shops/pkg"
	"shops/pkg/repository"
)

type Authorization interface {
	ConfirmUser(name, password string) (int, error)
	GenerateToken(id int) (string, error)
	ParseToken(token string) (int, error)
}

type Products interface {
	GetAllShops() ([]pkg.Shop, error)
	GetAllProducts() ([]pkg.Product, error)
	ReceiveProduct(prod pkg.Product, sc []pkg.ShopsProducts) (int, error)
	//UpdateProduct(prod *pkg.Product) error
}

type Receipts interface {
	AddToCart(userId int, cartItem *pkg.CartItem) error
	GetCarts(userId int) (*[]pkg.CartJSON, error)
	DeleteFromCart(item *pkg.CartItemsOnDeleteJSON, userId int) error
	CreateReceipt(shopId, userId, payOptId int) (int, error)
	TrySynchroByUserId(userId int) error
	GetReceipts(userId int) (*[]pkg.ReceiptJSON, error)
}

type Service struct {
	Authorization
	Products
	Receipts
	uConfs *UserServiceConfig
}

func InitNewService(uConfs *UserServiceConfig, repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(uConfs, repo),
		Products:      NewProductService(repo),
		Receipts:      NewReceiptsService(repo),
	}
}