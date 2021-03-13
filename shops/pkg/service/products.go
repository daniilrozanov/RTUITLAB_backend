package service

import (
	"encoding/json"
	"log"
	"shops/pkg"
	"shops/pkg/repository"
)


type ProductService struct {
	repo *repository.Repository
	rabbit *RabbitStruct
}

func NewProductService(repo *repository.Repository, rabbit *RabbitStruct) *ProductService {
	return &ProductService{repo: repo, rabbit: rabbit}
}

func (s *ProductService) StartConsume() error {
	q, err := s.rabbit.Channel.QueueDeclare("products", true, false, false, false, nil)
	if err != nil {
		return err
	}
	msgs, err := s.rabbit.Channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		log.Println("start product synchronization")
		for msg := range msgs {
			var cpd pkg.CreateProductData

			_ = json.Unmarshal(msg.Body, &cpd)
			if _, err := s.repo.ReceiveProduct(cpd.Prod, cpd.ShopsCount); err != nil {
				return
			}
		}
	}()
	return nil
}

func (p *ProductService) ReceiveProduct(prod pkg.Product, sc []pkg.ShopsProducts) (int, error) {
	return p.repo.ReceiveProduct(prod, sc)
}

func (p *ProductService) GetAllProducts() ([]pkg.Product, error) {
	return p.repo.GetAllProducts()
}

func (p *ProductService) GetAllShops() ([]pkg.Shop, error) {
	return p.repo.GetAllShops()
}
