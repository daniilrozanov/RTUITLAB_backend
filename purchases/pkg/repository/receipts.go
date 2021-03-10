package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	templates "purchases/pkg"
)

type ReceiptsService struct {
	db *sqlx.DB
}

func NewReceiptsService(db *sqlx.DB) *ReceiptsService {
	return &ReceiptsService{db: db}
}

func (r *ReceiptsService) CheckDBConnection() bool {
	log.Println("333")
	if r.db.Ping() != nil {
		return false
	}
	return true
}

func (r *ReceiptsService) InsertReceipt(urmap templates.UserReceiptMapJSON) error {
	var recId int

	query := fmt.Sprintf("INSERT INTO %s (body) VALUES ($1) RETURNING id", receiptsTable)
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	if err := tx.Get(&recId, query, urmap.Receipt); err != nil {
		tx.Rollback()
		return err
	}
	query = fmt.Sprintf("INSERT INTO %s (user_id, receipt_id) VALUES ($1, $2)", usersReceiptsTable)
	if _, err := tx.Exec(query, urmap.UserID, recId); err != nil {
		return err
	}
	tx.Commit()
	return nil
}
