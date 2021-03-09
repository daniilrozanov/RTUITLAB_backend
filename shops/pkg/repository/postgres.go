package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PostgresConfig struct {
	Host string
	Port string
	Username string
	Password string
	DBName string
	SSLMode string
}

const (
	shopsTable = "shops"
	productsTable = "products"
	receiptsTable = "receipts"
	shopsProductsTable = "shops_products"
	cartItemTable = "cart_item"
	cartsTable = "carts"
	reseiptsSynchroTable = "receipts_synchro"
	payOptionsTable = "pay_options"
)

func InitPostgresDB(cfg *PostgresConfig) (*sqlx.DB, error) {
	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password)
	db, err := sqlx.Open("postgres", dbUri)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
