package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shops/pkg"
)

func (r *ReceiptsService) AddToCart(userId int, cartItem *pkg.CartItem) error {
	var realQuantity, tmp int

	cartItem.UserID = userId
	if err := r.checkCompareDBQuantities(&realQuantity, cartItem); err != nil {
		return err
	}
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	if err := r.getCurrentCartNumberOrCreate(tx, cartItem); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := r.checkCartItemExists(tx, cartItem, &tmp); err == nil {
		cartItem.Quantity += tmp
		//log.Printf("finding existing cart_item...%d\n", cartItem.Quantity)
		if err := r.checkCompareQuantities(realQuantity, cartItem.Quantity); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := r.updateCartItem(tx, cartItem); err != nil {
			_ = tx.Rollback()
			return err
		}
		_ = tx.Commit()
		return nil
	} else if err == sql.ErrNoRows {
		//log.Printf("creating new cart_item...\n")
		if err := r.createCartItem(tx, cartItem); err != nil {
			_ = tx.Rollback()
			return err
		}
		_ = tx.Commit()
		return nil
	}
	_ = tx.Rollback()
	return err
	//проверки
}

func (r *ReceiptsService) createCartItem(tx *sqlx.Tx, cartItem *pkg.CartItem) error {
	var id int // to delete

	query := fmt.Sprintf("INSERT INTO %s (product_id, shop_id, user_id, quantity, index) VALUES ($1, $2, $3, $4, $5) RETURNING id", cartItemTable)
	err := tx.Get(&id, query, cartItem.ProductID, cartItem.ShopID, cartItem.UserID, cartItem.Quantity, cartItem.CartNumber)
	if err != nil {
		return err
	}
	//log.Printf("cart_item created: %d", id)
	return nil
}

func (r *ReceiptsService) updateCartItem(tx *sqlx.Tx, cartItem *pkg.CartItem) error {
	query := fmt.Sprintf("UPDATE %s SET quantity=$1 WHERE id=$2", cartItemTable)
	_, err := tx.Exec(query, cartItem.Quantity, cartItem.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReceiptsService) checkCartItemExists(tx *sqlx.Tx, cartItem *pkg.CartItem, tmp *int) error {
	query := fmt.Sprintf("SELECT id, quantity FROM %s WHERE product_id=$1 AND shop_id=$2 AND user_id=$3 AND index=$4", cartItemTable)
	*tmp = cartItem.Quantity
	err := tx.Get(cartItem, query, cartItem.ProductID, cartItem.ShopID, cartItem.UserID, cartItem.CartNumber)
	if err != nil {
		//log.Printf("cart_item doesnt exists\n")
		return err
	}
	//log.Printf("cart_item exists\n")
	return nil
}

func (r *ReceiptsService) checkCompareQuantities(realQuantity, quantity int) error {
	if quantity > realQuantity {
		return errors.New("not enough quantity of the product in the shop")
	}
	return nil
}

func (r *ReceiptsService) checkCompareDBQuantities(realQuantity *int, cartItem *pkg.CartItem) error {
	query := fmt.Sprintf("SELECT quantity FROM %s WHERE product_id=$1 AND shop_id=$2", shopsProductsTable)
	err := r.db.Get(realQuantity, query, cartItem.ProductID, cartItem.ShopID)
	if err == sql.ErrNoRows || *realQuantity < cartItem.Quantity {
		return errors.New("not enough quantity of the product in the shop")
	}
	return err
}

func (r *ReceiptsService) getCurrentCartNumberOrCreate(tx *sqlx.Tx, cartItem *pkg.CartItem) error {
	query := fmt.Sprintf("SELECT carts_number FROM %s WHERE user_id=$1 AND shop_id=$2 LIMIT 1", userCartsTable)
	err := tx.Get(&(cartItem.CartNumber), query, cartItem.UserID, cartItem.ShopID)
	if err == sql.ErrNoRows {
		query = fmt.Sprintf("INSERT INTO %s (user_id, shop_id, carts_number) VALUES ($1, $2, 1)", userCartsTable)
		_, err = tx.Exec(query, cartItem.UserID, cartItem.ShopID)
		if err != nil {
			return err
		}
		cartItem.CartNumber = 1
		//log.Printf("user_carts created: %d\n", cartItem.CartNumber)
		return nil
	}
	//log.Printf("user_carts founded: %d\n", cartItem.CartNumber)
	return err
}
