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
	AddToCart(int, *pkg.CartItemJSON) error
	GetCarts(int) (*[]pkg.CartJSON, error)
	DeleteFromCart(*pkg.CartItemsOnDeleteJSON, int) error
	CreateReceipt(int, int, int) (int, error)
	SetReceiptsSynchro(*[]int) error
	GetUserReceiptMap(*[]int) ([]pkg.UserReceiptMapJSON, error)
	GetReceipts(userId int) (*[]pkg.ReceiptJSON, error)
	GetUnsyncReceiptsIds(userId int) ([]int, error)
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
