package repository

import (
	templates "purchases/pkg"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type ProductPostgres struct {
	db *sqlx.DB
}

func NewProductPostgres(db *sqlx.DB) *ProductPostgres {
	return &ProductPostgres{db: db}
}

func (p *ProductPostgres) CreateProduct(userId int, prod *templates.Product) (int, error){

	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	createProdQuery := fmt.Sprintf("INSERT INTO %s (title, cost) values ($1, $2) RETURNING id", productsTable)
	row := tx.QueryRow(createProdQuery, prod.Title, prod.Cost) // execute sql and returns formatted sql values as response
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	createUserProdQuery := fmt.Sprintf("INSERT INTO %s (user_id, product_id) values ($1, $2)", usersProductsTable)
	_, err = tx.Exec(createUserProdQuery, userId, id) //execute sql query without checking response
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}

func (p *ProductPostgres) GetProductById(userId, prodId int) (templates.Product, error) {
	query := fmt.Sprintf(`SELECT pt.* FROM %s pt
	INNER JOIN %s upt ON pt.id = upt.product_id WHERE upt.user_id = $1 AND upt.product_id = $2`,
		productsTable, usersProductsTable)
	var out templates.Product
	err := p.db.Get(&out, query, userId, prodId)
	return out, err
}

func (p *ProductPostgres) GetAllProducts(userId int) ([]templates.Product, error){
	query := fmt.Sprintf(`SELECT pt.* FROM %s pt
	INNER JOIN %s upt ON pt.id = upt.product_id WHERE upt.user_id = $1`, productsTable, usersProductsTable)
	var out []templates.Product
	err := p.db.Select(&out, query, userId)
	return out, err
}

func (p *ProductPostgres) UpdateProduct(userId, prodId int, input *templates.UpdateProductInput) error{
	jParts := make([]string, 0)
	args := make([]interface{}, 0)
	i := 1

	if input.Title != nil {
		jParts = append(jParts, fmt.Sprintf("title=$%d", i))
		args = append(args, *input.Title)
		i++
	}
	if input.Cost != nil {
		jParts = append(jParts, fmt.Sprintf("cost=$%d", i))
		args = append(args, *input.Cost)
		i++
	}
	if input.Category != nil {
		jParts = append(jParts, fmt.Sprintf("category=$%d", i))
		args = append(args, *input.Category)
		i++
	}
	prequery := strings.Join(jParts, ", ")
	query := fmt.Sprintf("UPDATE %s pt SET %s FROM %s upt WHERE pt.id = upt.product_id AND upt.user_id = $%d AND upt.product_id = $%d",
		productsTable, prequery, usersProductsTable, i, i + 1)
	args = append(args, userId, prodId)
	_, err := p.db.Exec(query, args...)
	return err
}

func (p *ProductPostgres) DeleteProduct(userId, prodId int) error{
	query := fmt.Sprintf("DELETE FROM %s pt USING %s upt WHERE pt.id = upt.product_id AND pt.id = $1 AND upt.user_id = $2", productsTable, usersProductsTable)
	_, err := p.db.Exec(query, prodId, userId)
	return err
}