package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"shops/pkg"
	"time"
)

func (r *ReceiptsService) GetUserReceiptMap(recIds *[]int) (*[]pkg.UserReceiptMapJSON, error) {
	var urMap []pkg.UserReceiptMapJSON
	var userIds []int
	var recs *[]pkg.ReceiptJSON

	query := fmt.Sprintf("SELECT user_id FROM %s WHERE id IN (SELECT cart_id FROM %s WHERE id IN ($1))", cartsTable, receiptsTable)
	if err := r.db.Get(&userIds, query, *recIds); err != nil {
		return nil, err
	}
	recs, err := r.getReceiptsByIds(recIds)
	if err != nil {
		return nil, err
	}
	for i, _ := range *recs {
		urMap = append(urMap, pkg.UserReceiptMapJSON{
			UserID:  userIds[i],
		})
	}
	return &urMap, nil
}

func (r *ReceiptsService) GetReceipts(userId int)  (*[]pkg.ReceiptJSON, error) {
	var recIds []int

	query := fmt.Sprintf("SELECT id FROM %s WHERE cart_id IN (SELECT id FROM %s WHERE user_id=$1)", receiptsTable, cartsTable)
	if err := r.db.Select(&recIds, query, userId); err != nil {
		return nil, err
	}
	if len(recIds) == 0 {
		return nil, errors.New("user dont have any receipts")
	}
	log.Println("recids : ", recIds)
	return r.getReceiptsByIds(&recIds)
}

func (r *ReceiptsService) getReceiptsByIds(recIds *[]int) (*[]pkg.ReceiptJSON, error) {
	var recs []pkg.ReceiptJSON
	var carts *[]pkg.CartJSON
	var cartIds []int
	var times []time.Time
	var payOpts []string


	query := fmt.Sprintf("SELECT create_date FROM %s WHERE id IN (?)", receiptsTable)
	query, args, err := sqlx.In(query, *recIds)
	query = r.db.Rebind(query)
	if err != nil {
		return nil, err
	}
	if err := r.db.Select(&times, query, args...); err != nil {
		return nil, err
	}
	//log.Println("times : ", times)
	query = fmt.Sprintf("SELECT option FROM %s po JOIN %s rt ON rt.payopt_id=po.id WHERE rt.id IN (?)", payOptionsTable, receiptsTable)
	query, args, err = sqlx.In(query, *recIds)
	query = r.db.Rebind(query)
	if err != nil {
		return nil, err
	}
	if err := r.db.Select(&payOpts, query, args...); err != nil {
		return nil, err
	}
	//log.Println("payOpts : ", payOpts)
	query = fmt.Sprintf("SELECT cart_id FROM %s WHERE id IN (?)", receiptsTable)
	query, args, err = sqlx.In(query, *recIds)
	query = r.db.Rebind(query)
	if err != nil {
		return nil, err
	}
	if err := r.db.Select(&cartIds, query, args...); err != nil {
		return nil, err
	}
	//log.Println("times : ", times)
	carts, err = r.getCartsList(&cartIds)
	//log.Println("carts : ", *carts)
	if err != nil {
		return nil, err
	}
	log.Println(len(*carts), " ", len(cartIds), " ", len(times), " ", len(payOpts))
	for i, x := range *carts {
		recs = append(recs, pkg.ReceiptJSON{
			CartJSON: x,
			PayOption: payOpts[i],
			CreatedTime: times[i],
		})
		log.Println("N: ", i, "payopt: ", payOpts[i], "ceratedate: ", times[i], "cartjson: ", x)
	}
	return &recs, nil
}