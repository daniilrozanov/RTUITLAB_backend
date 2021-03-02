package service

import (
	"shops/pkg"
	"shops/pkg/repository"
)

type ProductService struct {
	repo *repository.Repository
}

func NewProductService(repo *repository.Repository) *ProductService {
	return &ProductService{repo: repo}
}

func (p *ProductService) ReceiveProduct(prod pkg.Product, sc []pkg.ShopsProducts) (int, error){
	return p.repo.ReceiveProduct(prod, sc)
}

func (p *ProductService) GetAllProducts() ([]pkg.Product, error) {
	return p.repo.GetAllProducts()
}

func (p *ProductService) GetAllShops() ([]pkg.Shop, error) {
	return p.repo.GetAllShops()
}


