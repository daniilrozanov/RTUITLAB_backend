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
	CreateReceipt()
	DeleteReceipt()
	GetReceips()
}

type Repository struct {
	Products
	Receipts
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Products: &ProductPostgres{db: db},
	}
}

