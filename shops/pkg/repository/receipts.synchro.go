package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func (r *ReceiptsService) SetReceiptsSynchro(recIds *[]int) error {
	query := fmt.Sprintf("UPDATE %s SET is_synchro=TRUE WHERE receipt_id IN (?)", reseiptsSynchroTable)
	query, args, err := sqlx.In(query, *recIds)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	if _, err := r.db.Exec(query, args...); err != nil {
		return err
	}
	return nil
}

func (r *ReceiptsService) GetUnsynchronizedReceiptsIds(userId int) (*[]int, error) {
	var recIds []int
	if userId < 1 {
		query := fmt.Sprintf("SELECT receipt_id FROM %s WHERE is_synchro=FALSE")
		if err := r.db.Get(&recIds, query); err != nil {
			return nil, err
		}
		return &recIds, nil
	}
	query := fmt.Sprintf("SELECT receipt_id FROM %s rc JOIN %s ON receipt_id=rc.id JOIN %s ct ON rc.cart_id=ct.id WHERE is_synchro=FALSE AND user_id=$1",
		receiptsTable, reseiptsSynchroTable, cartsTable)
	if err := r.db.Select(&recIds, query, userId); err != nil {
		return nil, err
	}
	return &recIds, nil
}