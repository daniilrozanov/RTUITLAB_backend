package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shops/pkg"
	"strconv"
	"strings"
)

type ProductPostgres struct {
	db *sqlx.DB
}

func (p *ProductPostgres) ReceiveProduct(prod pkg.Product, sc []pkg.ShopsProducts) (int, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}

	
	//ВСТАВКА ТОВАРА, ЕСЛИ ЕГО НЕТ В БАЗЕ ДАННЫХ
	if err := p.insertIfNotExists(tx, &prod); err != nil {
		tx.Rollback()
		return 0, err
	}
	values := make([]string, len(sc))
	for i := 0; i < len(sc); i++ {
		sc[i].ProductID = prod.ID
		values[i] = "( " + strconv.Itoa(sc[i].ProductID) + ", " + strconv.Itoa(sc[i].ShopID) + ", " + strconv.Itoa(sc[i].Quantity) + " )"
	}
	valuestr := strings.Join(values, ", ")
	query := fmt.Sprintf("INSERT INTO %s (product_id, shop_id, quantity) VALUES %s", shopsProductsTable, valuestr)
	if _, err = tx.Exec(query); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return prod.ID, nil
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