package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shops/pkg"
)

type ProductPostgres struct {
	db *sqlx.DB
}

func (p *ProductPostgres) ReceiveProduct(prod pkg.Product, sc []pkg.ShopsProducts) error {
	var pId int

	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("INSERT INTO %s (title, description, cost, category, code) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (code) DO UPDATE SET cost = EXCLUDED.cost RETURNING id", productsTable)
	if err := tx.Get(&pId, query, prod.Title, prod.Description, prod.Cost, prod.Category, prod.Code); err != nil {
		tx.Rollback()
		return errors.New("error while insert or update product: " + err.Error())
	}
	for _, x := range sc {
		query = fmt.Sprintf("INSERT INTO %s (product_id, shop_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (product_id, shop_id) DO UPDATE SET quantity=%s.quantity+$3", shopsProductsTable, shopsProductsTable)
		if _, err := tx.Exec(query, pId, x.ShopID, x.Quantity); err != nil {
			tx.Rollback()
			return errors.New("error while insert or update product quantities: " + err.Error())
		}
	}
	tx.Commit()
	return nil
}

func (p *ProductPostgres) GetAllProducts() ([]pkg.Product, error) {
	var products []pkg.Product

	query := fmt.Sprintf("SELECT * FROM %s", productsTable)

	if err := p.db.Select(&products, query); err != nil {
		return nil, err
	}
	for i, x := range products {
		query = fmt.Sprintf("SELECT shop_id, quantity FROM %s WHERE product_id=%d", shopsProductsTable, x.ID)
		if err := p.db.Select(&products[i].SQ, query); err != nil {
			return nil, err
		}
	}
	return products, nil
}

func (p *ProductPostgres) insertIfNotExists(tx *sql.Tx, prod *pkg.Product) error {
	query := fmt.Sprintf("SELECT id FROM %s WHERE code=%d", productsTable, prod.Code)
	err := tx.QueryRow(query).Scan(&(prod.ID))
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		query := fmt.Sprintf("INSERT INTO %s (title, description, cost, category, code) VALUES ($1, $2, $3, $4, $5) RETURNING id", productsTable)
		err = tx.QueryRow(query, prod.Title, prod.Description, prod.Cost, prod.Category, prod.Code).Scan(&(prod.ID))
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func (p *ProductPostgres) GetAllShops() ([]pkg.Shop, error) {
	var shops []pkg.Shop

	query := fmt.Sprintf("SELECT * FROM %s", shopsTable)
	if err := p.db.Select(&shops, query); err != nil {
		return nil, err
	}
	return shops, nil
}