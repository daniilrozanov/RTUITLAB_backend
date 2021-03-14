package repository

import (
	"errors"
	"fmt"
	"strings"
)

func (r *ReceiptsService) CreateReceipt(shopId, userId, payOptionId int) (int, error) { // смотреть количество и наличие товаров
	var ucId, recId, number int

	if err := r.checkPayOption(payOptionId); err != nil {
		return 0, err
	}
	if err := r.checkCart(&ucId, &number, shopId, userId); err != nil {
		return 0, err
	}
	if err := r.checkQuantity(ucId, shopId); err != nil {
		return 0, err
	}
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (cart_id, payopt_id) VALUES ($1, $2) RETURNING id", receiptsTable)
	if err := tx.Get(&recId, query, ucId, payOptionId); err != nil {
		tx.Rollback()
		return 0, errors.New("error while creating receipt: " + err.Error())
	}
	query = fmt.Sprintf("INSERT INTO %s (user_id, shop_id, number) VALUES ($1, $2, $3)", cartsTable)
	number += 1
	if _, err := tx.Exec(query, userId, shopId, number); err != nil {
		tx.Rollback()
		return 0, errors.New("error while creating new cart: " + err.Error())
	}
	var s int
	query = fmt.Sprintf("INSERT INTO %s (receipt_id, is_synchro) VALUES ($1, FALSE) RETURNING id", reseiptsSynchroTable)
	if err := tx.Get(&s, query, recId); err != nil {
		tx.Rollback()
		return 0, errors.New("error while creating synchro: " + err.Error())
	}
	query = fmt.Sprintf("UPDATE %s AS sp SET quantity=sp.quantity-ci.quantity FROM %s AS ci WHERE ci.cart_id=$1 AND ci.product_id=sp.product_id",
		shopsProductsTable, cartItemTable)
	if _, err := tx.Exec(query, ucId); err != nil {
		tx.Rollback()
		return 0, errors.New("error while updating products quantity in the shop: " + err.Error())
	}
	tx.Commit()
	return recId, nil
}

func (r *ReceiptsService) checkPayOption(payOptionId int) error {
	var x int

	query := fmt.Sprintf("SELECT id FROM %s WHERE id=$1", payOptionsTable)
	if err := r.db.Get(&x, query, payOptionId); err != nil {
		return errors.New("incorrect payment option")
	}
	return nil
}

func (r *ReceiptsService) checkCart(ucId, number *int, shopId, userId int) error {
	var x int

	query := fmt.Sprintf("SELECT id, number FROM %s WHERE shop_id=$1 AND user_id=$2 ORDER BY number DESC LIMIT 1", cartsTable)
	if err := r.db.QueryRow(query, shopId, userId).Scan(ucId, number); err != nil {
		return errors.New("cart doesn't exists") // если нет самой корзины
	}
	query = fmt.Sprintf("SELECT id FROM %s WHERE cart_id=$1 LIMIT 1", cartItemTable)
	if err := r.db.Get(&x, query, ucId); err != nil {
		return errors.New("cart is empty") // если в корзине нет товаров
	}
	return nil
}
func (r *ReceiptsService) checkQuantity(ucId, shopId int) error {
	var titles []string
	//var a, b int

	query := fmt.Sprintf("SELECT title FROM %s WHERE id IN (SELECT ci.product_id FROM %s ci JOIN %s sp ON ci.product_id=sp.product_id WHERE ci.quantity > sp.quantity AND cart_id=$1 AND sp.shop_id=$2)",
		productsTable, cartItemTable, shopsProductsTable)
	err := r.db.Select(&titles, query, ucId, shopId)
	if err != nil {
		return err
	}
	//query = fmt.Sprintf("SELECT ci.quantity, sp.quantity FROM %s ci JOIN %s sp ON ci.product_id=sp.product_id WHERE ci.quantity > sp.quantity AND cart_id=$1 AND sp.shop_id=$2", cartItemTable, shopsProductsTable)
	//r.db.QueryRow(query, ucId, shopId).Scan(&a, &b)
	//log.Println("ci.quantity, sp.quantity: ", a, " ", b)
	if len(titles) > 0 {
		return errors.New("quantities in " + strings.Join(titles, ", ") + " is bigger than available")
	}
	return nil
}
