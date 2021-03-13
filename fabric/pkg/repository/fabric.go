package repository

import (
	"database/sql"
	"fmt"
)

type Fabric struct {
	db *sql.DB
}

func NewFabric(db *sql.DB) *Fabric {
	return &Fabric{db: db}
}

func (f *Fabric) ProduceProduct(code, power int) error {
	query := fmt.Sprintf("INSERT INTO %s (code, powwer) VALUES ($1, $2) ON CONFLICT (code) DO UPDATE SET powwer = %s.powwer + $2", productsQuantityTable, productsQuantityTable)
	if _, err := f.db.Exec(query, code, power); err != nil {
		return err
	}
	return nil
}

func (f *Fabric) CompareQuantity (code, required int) (int, error) {
	query := fmt.Sprintf("SELECT powwer FROM %s WHERE code=$1", productsQuantityTable)
	var real int
	if err := f.db.QueryRow(query, code).Scan(&real); err != nil {
		return 0, err
	}
	if real >= required {
		return real - required, nil
	}
	return -1, nil
}
