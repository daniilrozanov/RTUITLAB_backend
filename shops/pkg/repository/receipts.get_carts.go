package repository

import (
	"fmt"
	"shops/pkg"
)

func (r *ReceiptsService) GetCarts(userId int) (*[]pkg.CartJSON, error) {
	var cartIds []int

	query := fmt.Sprintf("SELECT DISTINCT ON (shop_id) id FROM %s WHERE user_id=$1 ORDER BY shop_id, number DESC", cartsTable)
	err := r.db.Select(&cartIds, query, userId)
	if err != nil {
		return nil, err
	}
	return r.getCartsList(&cartIds)
}

func (r *ReceiptsService) getCartsList(cartIds *[]int) (*[]pkg.CartJSON, error) {
	var carts []pkg.CartJSON

	for _, cartId := range *cartIds {
		var cart pkg.CartJSON

		// Заполнение информации о магазине корзины
		query := fmt.Sprintf("SELECT * FROM %s WHERE id=(SELECT shop_id FROM %s WHERE id=$1)", shopsTable, cartsTable)
		if err := r.db.Get(&(cart.Shop), query, cartId); err != nil {
			return nil, err
		}
		// Получение списка объектов данной корзины
		query = fmt.Sprintf("SELECT product_id, quantity FROM %s WHERE cart_id=$1", cartItemTable)
		rowssh, err := r.db.Queryx(query, cartId)
		if err != nil {
			return nil, err
		}
		//Проход по списку объектов данной корзины
		for rowssh.Next() {
			var cartItem pkg.CartItem
			var prod pkg.ProductJSON

			err := rowssh.StructScan(&cartItem)
			if err != nil {
				return nil, err
			}
			//Получение информации о товаре и его количестве
			query := fmt.Sprintf("SELECT p.id, title, cost, coalesce(cp.category, p.category) FROM %s p LEFT JOIN %s cp ON p.id=cp.product_id WHERE p.id=$1 AND (cp.cart_id=$2 OR cp.cart_id IS NULL)", productsTable, productsCustomCategoriesTable)
			if err := r.db.QueryRow(query, cartItem.ProductID, cartId).Scan(&prod.ID, &prod.Title, &prod.Cost, &prod.Category); err != nil {
				return nil, err
			}
			//query := fmt.Sprintf("SELECT category FROM %s WHERE cart_id=$1 AND product_id=$2", productsCustomCategoriesTable)
			//if err := r.db.Get(&prod, query)
			prod.Quantity = cartItem.Quantity
			prod.EntireCost = prod.Cost * cartItem.Quantity
			cart.SummaryCost += prod.EntireCost
			cart.Products = append(cart.Products, prod)

		}
		if err := rowssh.Err(); err != nil {
			return nil, err
		}
		if len(cart.Products) > 0 {
			carts = append(carts, cart)
		}
	}
	return &carts, nil
}