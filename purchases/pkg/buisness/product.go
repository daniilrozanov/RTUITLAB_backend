package buisness

import (
	templates "purchases/pkg"
	"purchases/pkg/repository"
)

type ProductService struct {
	repo repository.ProductLogging
}

func NewProductService(repo repository.ProductLogging) *ProductService {
	return &ProductService{repo: repo}
}

func (p *ProductService) CreateProduct(userId int, prod *templates.Product) (int, error) {
	return p.repo.CreateProduct(userId, prod)
}

func (p *ProductService) GetProductById(userId, prodId int) (templates.Product, error) {
	return p.repo.GetProductById(userId, prodId)
}

func (p *ProductService) GetAllProducts(userId int) ([]templates.Product, error) {
	return p.repo.GetAllProducts(userId)
}

func (p *ProductService) UpdateProduct(userId, prodId int, input *templates.UpdateProductInput) error {
	return p.repo.UpdateProduct(userId, prodId, input)
}

func (p *ProductService) DeleteProduct(userId, prodId int) error {
	return p.repo.DeleteProduct(userId, prodId)
}