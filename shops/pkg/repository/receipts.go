package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shops/pkg"
)

type ReceiptsService struct {
	db *sqlx.DB
}

func (r *ReceiptsService) GetCarts(userId int) ([]pkg.CartJSON, error){
	var cartItems *[][]pkg.CartItem

	cartItems, err := r.getCartItemsByUserId(userId)
	if err != nil {
		return nil, err
	}
	return r.processCartItemsToJSON(cartItems)
}

func (r *ReceiptsService) processCartItemsToJSON(cartItems *[][]pkg.CartItem) ([]pkg.CartJSON, error){
	var carts []pkg.CartJSON

	for _, x := range *cartItems {
		var cart pkg.CartJSON
		var shop pkg.Shop

		query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", shopsTable)
		if err := r.db.Get(&shop, query, x[0].ShopID); err != nil {
			return nil, err
		}
		cart.Shop = shop
		cart.Index = x[0].CartNumber
		for _, y := range x {
			var prod pkg.ProductJSON

			query := fmt.Sprintf("SELECT pr.id, title, cost, category, quantity FROM %s pr JOIN %s ci ON ci.product_id=pr.id WHERE pr.id=$1 AND shop_id=$2 AND user_id=$3 AND index=$4", productsTable, cartItemTable)
			if err := r.db.Get(&prod, query, y.ProductID, y.ShopID, y.UserID, y.CartNumber); err != nil {
				return nil, err
			}
			cart.Products = append(cart.Products, prod)
			cart.SummaryCost += prod.Cost * prod.Quantity
		}
		carts = append(carts, cart)
	}
	return carts, nil
}

func (r *ReceiptsService) getCartItemsByUserId(userId int) (*[][]pkg.CartItem, error){
	var ret [][]pkg.CartItem

	query := fmt.Sprintf("SELECT shop_id, carts_number FROM %s WHERE user_id=$1", userCartsTable)
	rowsal, err := r.db.Queryx(query, userId) // пары магазин-корзина пользователя. магазины не повторяются
	if err != nil {
		return nil, err
	}
	for rowsal.Next() { // проход по каждому магазину-корзине
		var userCart pkg.UserCarts // текущий магазин-корзина
		var cartItems []pkg.CartItem // объекты корзины для текущего магазина

		rowsal.StructScan(&userCart)
		query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1 AND shop_id=$2 AND index=$3", cartItemTable)
		rowssh, err := r.db.Queryx(query, userId, userCart.ShopID, userCart.NumberOfCarts) // sql объекты корзины пользователя для текущего магазина
		if err != nil {
			return nil, err
		}
		for rowssh.Next() { // проход по объектам корзины пользователя текущего магазина
			var cartItem pkg.CartItem //экземпляр объекта

			rowssh.StructScan(&cartItem)
			cartItems = append(cartItems, cartItem) // присоединение объекта корзины в список о-к для текущего магазина
		}
		if err := rowssh.Err(); err != nil {
			return nil, err
		}
		if len(cartItems) > 0 {
			ret = append(ret, cartItems) // присоединение объектов корзины для текущего магазина к списку всех объектов всех корзин пользователя
		}
	}
	if err := rowsal.Err(); err != nil {
		return nil, err
	}
	return &ret, nil
}

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
	_, err = r.db.Exec(query, realQuantity - item.Quantity, itemId)
	return err
}