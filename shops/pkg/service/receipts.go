package service

import (
	"shops/pkg"
	"shops/pkg/repository"
)

type ReceiptsService struct {
	repo *repository.Repository
}

func NewReceiptsService(repo *repository.Repository) *ReceiptsService {
	return &ReceiptsService{repo: repo}
}

func (r *ReceiptsService) AddToCart(userId int, cartItem *pkg.CartItem) error{
	if cartItem.Quantity == 0 {
		cartItem.Quantity = 1
	}
	return r.repo.AddToCart(userId, cartItem)
}

func (r *ReceiptsService) GetCarts(userId int) ([]pkg.CartJSON, error) {
	return r.repo.GetCarts(userId)
}

func (r *ReceiptsService) DeleteFromCart (item *pkg.CartItemsOnDeleteJSON, userID int) error{
	return r.repo.DeleteFromCart(item, userID)
}

