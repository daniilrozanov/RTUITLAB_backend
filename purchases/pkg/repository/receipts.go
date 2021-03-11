package repository

import (
	"encoding/json"
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
	var data []byte

	data, err := json.Marshal(urmap.Receipt)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("INSERT INTO %s (body) VALUES ($1) RETURNING id", receiptsTable)
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	if err := tx.Get(&recId, query, data); err != nil {
		tx.Rollback()
		return err
	}
	query = fmt.Sprintf("INSERT INTO %s (user_id, receipt_id) VALUES ($1, $2)", usersReceiptsTable)
	if _, err := tx.Exec(query, urmap.UserID, recId); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *ReceiptsService) GetReceipts (userID int) ([]templates.ReceiptJSON, error){
	var recsdb []templates.ReceiptDB
	var recs []templates.ReceiptJSON

	query := fmt.Sprintf("SELECT body FROM %s WHERE id IN (SELECT receipt_id FROM %s WHERE user_id=$1)", receiptsTable, usersReceiptsTable)
	if err := r.db.Select(&recsdb, query, userID); err != nil {
		return nil, err
	}

	for _, x := range recsdb {
		var rec templates.ReceiptJSON
		json.Unmarshal(x.Receipt, &rec)
		recs = append(recs, rec)
	}
	return recs, nil
}
