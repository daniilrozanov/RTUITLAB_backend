package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"shops/pkg"
	"shops/pkg/repository"
)

type ReceiptsService struct {
	uConfs *UserServiceConfig
	repo   *repository.Repository
}

func NewReceiptsService(repo *repository.Repository) *ReceiptsService {
	return &ReceiptsService{repo: repo}
}

func (r *ReceiptsService) AddToCart(userId int, cartItem *pkg.CartItem) error {
	if cartItem.Quantity == 0 {
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

func (r *ReceiptsService) GetReceipts(userId int) (*[]pkg.ReceiptJSON, error) {
	return r.repo.GetReceipts(userId)
}

func (r *ReceiptsService) TrySynchroByUserId(userId int) error {
	unsyncRecIds, err := r.repo.GetUnsynchronizedReceiptsIds(userId)
	if err != nil {
		return err
	}
	receipts, err := r.repo.GetUserReceiptMap(unsyncRecIds)
	if err != nil {
		return err
	}
	byteJSON, err := json.Marshal(*receipts)
	if err != nil {
		return err
	}
	cryptedJSON, err := encrypt(byteJSON, UsersTransportKey)
	if err != nil {
		return err
	}
	response, err := http.Post(r.getUserServiceURI(), "application/json", bytes.NewBuffer(cryptedJSON))
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	decryptedBody, err := decrypt(b, UsersTransportKey)
	if err != nil {
		return err
	}
	var responseJSON map[string]interface{}
	err = json.Unmarshal(decryptedBody, responseJSON)
	if err != nil {
		return err
	}
	if responseJSON["status"] == "ok" {
		return nil
	}
	return errors.New("service users not available")
}

func (a *ReceiptsService) getUserServiceURI() string {
	if a.uConfs.Port == ":" {
		return a.uConfs.Scheme + a.uConfs.Host + "/" + a.uConfs.SynchroURN
	}
	return a.uConfs.Scheme + a.uConfs.Host + ":" + a.uConfs.Port + "/" + a.uConfs.SynchroURN
}
