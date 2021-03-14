package service

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"log"
	"shops/pkg"
	"shops/pkg/repository"
)

type RabbitStruct struct {
	Channel *amqp.Channel
}

type ReceiptsService struct {
	uConfs *UserServiceConfig
	rabbit *RabbitStruct
	repo   *repository.Repository
}

func NewReceiptsService(repo *repository.Repository, rabbitStruct *RabbitStruct, uConfs *UserServiceConfig) *ReceiptsService {
	return &ReceiptsService{repo: repo, rabbit: rabbitStruct, uConfs: uConfs}
}

func (r *ReceiptsService) AddToCart(userId int, cartItem *pkg.CartItemJSON) error {
	if cartItem.Quantity < 1 {
		cartItem.Quantity = 1
	}
	return r.repo.AddToCart(userId, cartItem)
}

func (r *ReceiptsService) GetCarts(userId int) (*[]pkg.CartJSON, error) {
	return r.repo.GetCarts(userId)
}

func (r *ReceiptsService) DeleteFromCart(item *pkg.CartItemsOnDeleteJSON, userID int) error {
	return r.repo.DeleteFromCart(item, userID)
}

func (r *ReceiptsService) CreateReceipt(shopId, userId, payOptId int) (int, error) {
	return r.repo.CreateReceipt(shopId, userId, payOptId)
}

func (r *ReceiptsService) SendReceiptToRabbit(recId int) error {
	var urMap []pkg.UserReceiptMapJSON
	recIds := []int{recId}
	urMap, err := r.repo.GetUserReceiptMap(&recIds)
	if err != nil {
		return err
	}
	data, err := json.Marshal(urMap[0])
	if err != nil {
		return err
	}

	q, err := r.rabbit.Channel.QueueDeclare("receipts", true, false, false, false, nil)
	err = r.rabbit.Channel.Publish("", q.Name, false, false, amqp.Publishing{ContentType: "application/json", Body: data})
	if err != nil {
		return errors.New("Failed to publish a message: " + err.Error())
	}
	return nil
}

func (r *ReceiptsService) SendUnsyncReceiptsToRabbit() error {
	recIds, err := r.repo.GetUnsyncReceiptsIds(0)
	log.Println(recIds)
	if err != nil {
		return err
	}
	for _, x := range recIds {
		if err := r.SendReceiptToRabbit(x); err != nil {
			return err
		}
	}
	if err := r.repo.SetReceiptsSynchro(&recIds); err != nil {
		return err
	}
	return nil
}

func (r *ReceiptsService) GetReceipts(userId int) (*[]pkg.ReceiptJSON, error) {
	return r.repo.GetReceipts(userId)
}

func (a *ReceiptsService) getUserServiceURI() string {
	if a.uConfs.Port == ":" {
		return a.uConfs.Scheme + a.uConfs.Host + "/" + a.uConfs.SynchroURN
	}
	return a.uConfs.Scheme + a.uConfs.Host + ":" + a.uConfs.Port + "/" + a.uConfs.SynchroURN
}

func (r *ReceiptsService) SetReceiptsSynchro(recIds *[]int) error {
	return r.repo.SetReceiptsSynchro(recIds)
}