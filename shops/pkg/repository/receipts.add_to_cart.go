package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shops/pkg"
)

func (r *ReceiptsService) AddToCart(userId int, cartItemJ *pkg.CartItemJSON) error {
	var realQuantity, tmp, cartId int
	cartItem := &cartItemJ.CartItem

	cartItem.UserID = userId
	if err := r.checkCompareDBQuantities(&realQuantity, cartItem); err != nil {
		return err
	}
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	if cartId, err = r.getCurrentCartNumberOrCreate(tx, cartItem); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := r.checkCartItemExists(tx, cartItem, &tmp); err == nil {
		cartItem.Quantity += tmp
		if err := r.checkCompareQuantities(realQuantity, cartItem.Quantity); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := r.updateCartItem(tx, cartItem); err != nil {
			_ = tx.Rollback()
			return err
		}
	} else if err == sql.ErrNoRows {
		if err := r.createCartItem(tx, cartItem, cartId); err != nil {
			_ = tx.Rollback()
			return err
		}
	} else {
		_ = tx.Rollback()
		return err
	}
	if cartItemJ.Category != "" {
		query := fmt.Sprintf("INSERT INTO %s (cart_id, product_id, category) VALUES ($1, $2, $3) ON CONFLICT (cart_id, product_id) DO UPDATE SET category=$3", productsCustomCategoriesTable)
		if _, err := r.db.Exec(query, cartId, cartItem.ProductID, cartItemJ.Category); err != nil {
			return err
		}
	}
	_ = tx.Commit()
	return nil
}

func (r *ReceiptsService) createCartItem(tx *sqlx.Tx, cartItem *pkg.CartItem, cartId int) error {
	query := fmt.Sprintf("INSERT INTO %s (product_id, quantity, cart_id) VALUES ($1, $2, $3) RETURNING id", cartItemTable)
	_, err := tx.Exec(query, cartItem.ProductID, cartItem.Quantity, cartId)
	if err != nil {
		return err
	}
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
	query := fmt.Sprintf("SELECT id, quantity FROM %s WHERE product_id=$1 AND cart_id=(SELECT id FROM %s WHERE shop_id=$2 AND user_id=$3 AND number=$4 LIMIT 1)", cartItemTable, cartsTable)
	*tmp = cartItem.Quantity
	err := tx.Get(cartItem, query, cartItem.ProductID, cartItem.ShopID, cartItem.UserID, cartItem.Number)
	if err != nil {
		//log.Printf("cart_item doesnt exists\n")
		return err
	}
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

func (r *ReceiptsService) getCurrentCartNumberOrCreate(tx *sqlx.Tx, cartItem *pkg.CartItem) (int, error) {
	var cartId int

	query := fmt.Sprintf("SELECT id, number FROM %s WHERE user_id=$1 AND shop_id=$2 ORDER BY number DESC", cartsTable)
	err := tx.QueryRowx(query, cartItem.UserID, cartItem.ShopID).Scan(&cartId, &cartItem.Number)
	if err == sql.ErrNoRows {
		query = fmt.Sprintf("INSERT INTO %s (user_id, shop_id, number) VALUES ($1, $2, 1) RETURNING id", cartsTable)
		err = tx.Get(&cartId, query, cartItem.UserID, cartItem.ShopID)
		if err != nil {
			return 0, err
		}
		cartItem.Number = 1
		return cartId, nil
	}
	return cartId, err
}
