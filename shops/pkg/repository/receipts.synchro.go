package repository

import (
	"fmt"
	"shops/pkg"
)

func (r *ReceiptsService) SetReceiptsSynchro(recIds []int) error {
	query := fmt.Sprintf("UPDATE %s SET is_synchro=TRUE WHERE receipt_id IN ($1)", userCartsTable)
	if _, err := r.db.Exec(query, recIds); err != nil {
		return err
	}
	return nil
}

func (r *ReceiptsService) GetUnsynchronizedReceiptsIds(userId int) ([]int, error) {
	var recIds []int
	if userId < 1 {
		query := fmt.Sprintf("SELECT receipt_id FROM %s WHERE is_synchro=FALSE")
		if err := r.db.Get(&recIds, query); err != nil {
			return nil, err
		}
		return recIds, nil
	}
	query := fmt.Sprintf("SELECT receipt_id FROM %s rc JOIN %s ON receipt_id=rc.id WHERE is_synchro=FALSE AND user_id=$1", receiptsTable, reseiptsSynchroTable)
	if err := r.db.Get(&recIds, query, userId); err != nil {
		return nil, err
	}
	return recIds, nil
}


func (r *ReceiptsService) GetUserReceiptMap(recIds []int) (*[]pkg.UserReceiptMapJSON, error) {
	var receipts []pkg.Receipt
	var urMap []pkg.UserReceiptMapJSON

	query := fmt.Sprintf("SELECT * FROM %s WHERE id IN ($1)")
	if err := r.db.Get(&receipts, query, recIds); err != nil {
		return nil, err
	}
	for _, x := range receipts {
		var urPair pkg.UserReceiptMapJSON

		query = fmt.Sprintf("SELECT user_id FROM %s WHERE id=$1", receiptsTable)
		if err := r.db.Get(&urPair.UserID, query, x.UserID); err != nil {
			return nil, err
		}
		query = fmt.Sprintf("SELECT * FROM %s WHERE id=(SELECT shop_id FROM %s WHERE id=$1 LIMIT 1)", shopsTable, receiptsTable)
		if err := r.db.Get(&urPair.Receipt.Shop, query, x.Id); err != nil {
			return nil, err
		}
		query = fmt.Sprintf("SELECT * FROM %s WHERE id IN (SELECT product_id FROM %s WHERE )", cartItemTable)
	}

}