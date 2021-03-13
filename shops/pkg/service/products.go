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

func (p *ProductService) StartConsume() error {
	q, err := p.rabbit.Channel.QueueDeclare("products", true, false, false, false, nil)
	if err != nil {
		return err
	}
	msgs, err := p.rabbit.Channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		log.Println("start product synchronization")
		for msg := range msgs {
			var cpd pkg.CreateProductData

			if err := json.Unmarshal(msg.Body, &cpd); err != nil {
				log.Println("error while parsing json: ", err.Error())
				return
			}
			if err := p.repo.ReceiveProduct(cpd.Prod, cpd.ShopsCount); err != nil {
				log.Println("synchronization with fabric stopped: ", err.Error())
				return
			}
			log.Println("received: ", cpd.Prod.Title)
		}
		log.Println("synchronization with fabric stopped.")
	}()
	return nil
}

func (p *ProductService) GetAllProducts() ([]pkg.Product, error) {
	return p.repo.GetAllProducts()
}

func (p *ProductService) GetAllShops() ([]pkg.Shop, error) {
	return p.repo.GetAllShops()
}
