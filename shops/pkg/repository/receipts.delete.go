package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"shops/pkg"
)

type ReceiptsService struct {
	db *sqlx.DB
}

func (r *ReceiptsService) DeleteFromCart(item *pkg.CartItemsOnDeleteJSON, userID int) error {
	var cartId, itemId, realQuantity int

	if item.Quantity < 0 {
		return errors.New("entered quantity to delete less than zero")
	}
	query := fmt.Sprintf("SELECT id FROM %s WHERE user_id=$1 AND shop_id=$2 ORDER BY number DESC LIMIT 1", cartsTable)
	log.Println("userId: ", userID, " shopId: ", item.ShopID)
	err := r.db.Get(&cartId, query, userID, item.ShopID)
	if err != nil {
		return errors.New("cart doesn't exists")
	}
	if item.Quantity == 0 {
		query := fmt.Sprintf("DELETE FROM %s WHERE cart_id=$1 AND product_id=$2", cartItemTable)
		_, err := r.db.Exec(query, cartId, item.ProductID)
		if err != nil {
			return errors.New("impossible to delete cart item")
		}
		return nil
	}
	query = fmt.Sprintf("SELECT id, quantity FROM %s WHERE product_id=$1 AND cart_id=$2", cartItemTable)
	if err := r.db.QueryRowx(query, item.ProductID, cartId).Scan(&itemId, &realQuantity); err != nil {
		return err
	}
	if realQuantity <= item.Quantity {
		query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", cartItemTable)
		_, err := r.db.Exec(query, itemId)
		if err != nil {
			return err
		}
		return nil
	}
	query = fmt.Sprintf("UPDATE %s SET quantity=$1 WHERE id=$2", cartItemTable)
	_, err = r.db.Exec(query, realQuantity-item.Quantity, itemId)
	return err
}