package repository

import (
	"fmt"
	"shops/pkg"
)

func (r *ReceiptsService) GetCarts(userId int) ([]pkg.CartJSON, error) {
	var cartItems *[][]pkg.CartItem

	cartItems, err := r.getCartItemsByUserId(userId)
	if err != nil {
		return nil, err
	}
	return r.cartItemsToCartsJSON(cartItems)
}

func (r *ReceiptsService) cartItemsToCartsJSON(cartItems *[][]pkg.CartItem) ([]pkg.CartJSON, error) {
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

func (r *ReceiptsService) getCartItemsByUserId(userId int) (*[][]pkg.CartItem, error) {
	var ret [][]pkg.CartItem

	query := fmt.Sprintf("SELECT shop_id, carts_number FROM %s WHERE user_id=$1", userCartsTable)
	rowsal, err := r.db.Queryx(query, userId) // пары магазин-корзина пользователя. магазины не повторяются
	if err != nil {
		return nil, err
	}
	for rowsal.Next() { // проход по каждому магазину-корзине
		var userCart pkg.UserCarts   // текущий магазин-корзина
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

