package repository

import (
	"github.com/jmoiron/sqlx"
	"shops/pkg"
)

type Products interface {
	GetAllShops() ([]pkg.Shop, error)
	GetAllProducts() ([]pkg.Product, error)
	ReceiveProduct(pkg.Product, []pkg.ShopsProducts) (int, error)
}

type Receipts interface {
	AddToCart(int, *pkg.CartItem) error
	GetCarts(int) ([]pkg.CartJSON, error)
	DeleteFromCart(item *pkg.CartItemsOnDeleteJSON, userID int) error
	//CreateReceipt(int) error
}

type Repository struct {
	Products
	Receipts
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Products: &ProductPostgres{db: db},
		Receipts: &ReceiptsService{db: db},
	}
}

