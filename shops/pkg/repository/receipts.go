package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shops/pkg"
)

type ReceiptsService struct {
	db *sqlx.DB
}


func (r *ReceiptsService) DeleteFromCart(item *pkg.CartItemsOnDeleteJSON, userID int) error {
	var currentIndex, itemId, realQuantity int

	query := fmt.Sprintf("SELECT carts_number FROM %s WHERE user_id=$1 AND shop_id=$2", userCartsTable)
	err := r.db.Get(&currentIndex, query, userID, item.ShopID)
	if err != nil {
		return err
	}
	if item.Quantity < 0 {
		return errors.New("entered quantity to delete less than zero")
	}
	if item.Quantity == 0 {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1 AND shop_id=$2 AND product_id=$3 AND index=$4", cartItemTable)
		_, err := r.db.Exec(query, userID, item.ShopID, item.ProductID, currentIndex)
		if err != nil {
			return err
		}
		return nil
	}
	query = fmt.Sprintf("SELECT id, quantity FROM %s WHERE user_id=$1 AND shop_id=$2 AND product_id=$3 AND index=$4", cartItemTable)
	if err := r.db.QueryRowx(query, userID, item.ShopID, item.ProductID, currentIndex).Scan(&itemId, &realQuantity); err != nil {
		return err
	}
	if realQuantity <= item.Quantity {
		query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", cartItemTable)
		_, err := r.db.Exec(query, userID, itemId)
		if err != nil {
			return err
		}
		return nil
	}
	query = fmt.Sprintf("UPDATE %s SET quantity=$1 WHERE id=$2", cartItemTable)
	_, err = r.db.Exec(query, realQuantity-item.Quantity, itemId)
	return err
}

func (r *ReceiptsService) CreateReceipt(shopId, userId int) (int, error) {
	var ucId, currentIndex, recId int

	query := fmt.Sprintf("SELECT id, carts_number FROM %s WHERE shop_id=$1 AND user_id=$2", userCartsTable)
	if err := r.db.QueryRowx(query, shopId, userId).Scan(&ucId, &currentIndex); err != nil {
		return 0, err
	}
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	query = fmt.Sprintf("INSERT INTO %s (user_id, shop_id, cart_item_number) VALUES ($1, $2, $3) RETURNING id", receiptsTable)
	if err := tx.Get(&recId, query, userId, shopId, currentIndex); err != nil {
		tx.Rollback()
		return 0, err
	}
	query = fmt.Sprintf("UPDATE %s SET carts_number=carts_number+1 WHERE id=$1", userCartsTable)
	if _, err := tx.Exec(query, ucId); err != nil {
		tx.Rollback()
		return 0, err
	}
	query = fmt.Sprintf("INSERT INTO %s (receipt_id, is_synchro) VALUES ($1, FALSE)")
	if _, err := tx.Exec(query, recId); err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return recId, err
}

func (r *ReceiptsService) cartItemsToReceiptJSON(cartItems *[]pkg.CartItem) {

}
