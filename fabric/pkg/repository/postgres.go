package repository

import (
	"database/sql"
	"fmt"
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
	productsQuantityTable = "products_quantity"
)

func InitPostgresDB(cfg *PostgresConfig) (*sql.DB, error) {
	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password)
	db, err := sql.Open("postgres", dbUri)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS products_quantity (id SERIAL PRIMARY KEY, code INTEGER UNIQUE, powwer INTEGER);")
	if _, err := db.Exec(query); err != nil {
		return db, err
	}
	return db, nil
}
