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
	AddToCart(userId int, cartItem *pkg.CartItemJSON) error
	GetCarts(userId int) (*[]pkg.CartJSON, error)
	DeleteFromCart(item *pkg.CartItemsOnDeleteJSON, userId int) error
	CreateReceipt(shopId, userId, payOptId int) (int, error)
	GetReceipts(userId int) (*[]pkg.ReceiptJSON, error)
	SendReceiptToRabbit(recId int) error
	SetReceiptsSynchro(*[]int) error
	SendUnsyncReceiptsToRabbit() error
}

type Service struct {
	Authorization
	Products
	Receipts
}

func InitNewService(uConfs *UserServiceConfig, repo *repository.Repository, rabbitStruct *RabbitStruct) *Service {
	return &Service{
		Authorization: NewAuthService(uConfs, repo),
		Products:      NewProductService(repo, rabbitStruct),
		Receipts:      NewReceiptsService(repo, rabbitStruct, uConfs),
	}
}