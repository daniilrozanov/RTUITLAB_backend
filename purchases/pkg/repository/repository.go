package repository

import (
	templates "purchases/pkg"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	GetUser(username, password string) (templates.User, error)
	CreateUser(user templates.User) (int, error)
}

type ProductLogging interface {
	CreateProduct(userId int, prod *templates.Product) (int, error)
	GetProductById(userId, prodId int) (templates.Product, error)
	GetAllProducts(userId int) ([]templates.Product, error)
	UpdateProduct(userId, prodId int, input *templates.UpdateProductInput) error
	DeleteProduct(userId, prodId int) error
}

type Receipts interface {
	CheckDBConnection() bool
	InsertReceipt(urmap templates.UserReceiptMapJSON) error
	GetReceipts (userID int) ([]templates.ReceiptJSON, error)
}

type Repository struct {
	Authorization
	ProductLogging
	Receipts
}

func InitRepositoryLayer (db *sqlx.DB) *Repository{
	return &Repository{
		Authorization: NewAuthPostgres(db),
		ProductLogging: NewProductPostgres(db),
		Receipts: NewReceiptsService(db),
	}
}
