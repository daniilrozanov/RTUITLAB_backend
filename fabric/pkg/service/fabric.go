package service

import (
	"encoding/json"
	"errors"
	"fabric/pkg"
	"fabric/pkg/repository"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"time"
)

type RabbitStruct struct {
	Channel *amqp.Channel
}

type Fabric struct {
	repo *repository.Fabric
	rabbit *RabbitStruct
	cfg *ConfigJSON
}

func NewFabric(repo *repository.Fabric, rabbit *RabbitStruct, cfg *ConfigJSON) *Fabric {
	return &Fabric{repo: repo, rabbit: rabbit, cfg: cfg}
}

type ConfigJSON struct {
	Products []pkg.FabricProductJSON `json:"products"`
	Period int
}

func InitConfig(filename string) (*ConfigJSON, error) {
	plan, _ := ioutil.ReadFile(filename)
	var cfg ConfigJSON
	if err := json.Unmarshal(plan, &cfg); err != nil {
		return nil, err
	}
	 return &cfg, nil
}

func (f *Fabric) ProduceProduct(code, power int) error {
	return f.repo.ProduceProduct(code, power)
}

func (f *Fabric) CompareQuantity (code, required int) (int, error) {
	return f.repo.CompareQuantity(code, required)
}

func (f *Fabric) SendToRabbit(cpd pkg.CreateProductData) error {
	data, err := json.Marshal(cpd)
	if err != nil {
		return err
	}
	if f.rabbit.Channel == nil {
		return errors.New("no rabbitmq connection")
	}
	q, err := f.rabbit.Channel.QueueDeclare("products", true, false, false, false, nil)
	if err != nil {
		return errors.New("Failed to connect queue: " + err.Error())
	}
	err = f.rabbit.Channel.Publish("", q.Name, false, false, amqp.Publishing{ContentType: "application/json", Body: data})
	if err != nil {
		return errors.New("Failed to publish a message: " + err.Error())
	}
	return nil
}

func (f *Fabric) getRequiredQuantity(sc []pkg.ShopsProducts) int {
	var res int
	for _, x := range sc {
		res += x.Quantity
	}
	return res
}

func (f *Fabric) StartProducing() error {
	for {
		log.Println("")
		for _, x := range f.cfg.Products {
			if err := f.ProduceProduct(x.CreateProductData.Prod.Code, x.Power); err != nil {
				log.Println("product processing error: ", err.Error())
				return err
			}
			log.Printf("product: %-40s | produced: %5d\n", x.Prod.Title, x.Power)
			var err error
			var diff, reqq int
			reqq = f.getRequiredQuantity(x.ShopsCount)
			if diff, err = f.CompareQuantity(x.CreateProductData.Prod.Code, reqq); err != nil {
				log.Println("product processing error: ", err.Error())
				return err
			}
			if diff > -1 {
				if err := f.ProduceProduct(x.CreateProductData.Prod.Code, -reqq); err != nil {
					log.Println("product processing error: ", err.Error())
					return err
				}
			} else {
				continue
			}
			if err := f.SendToRabbit(x.CreateProductData); err != nil {
				if err := f.ProduceProduct(x.CreateProductData.Prod.Code, reqq); err != nil {
					log.Println("product processing error: ", err.Error())
					return err
				}
				continue
			}
			log.Printf("product: %-40s | sent: %5d\n", x.Prod.Title, reqq)
		}
		time.Sleep(time.Duration(f.cfg.Period) * time.Second)
	}
}
