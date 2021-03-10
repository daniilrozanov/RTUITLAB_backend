package buisness

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	templates "purchases/pkg"
	"purchases/pkg/repository"
	"time"
)

type SynchroService struct {
	rabbit RabbitStruct
	repo   repository.Receipts
}

func NewSynchroService(rabbit RabbitStruct, repo repository.Repository) *SynchroService {
	return &SynchroService{rabbit: rabbit, repo: repo}
}

type RabbitStruct struct {
	Channel *amqp.Channel
}

type RabbitConnectionConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewRabbitStruct(confs *RabbitConnectionConfig) (RabbitStruct, error) {
	var rabbit RabbitStruct

	conn, err := amqp.Dial("amqp://" + confs.Username + ":" + confs.Password + "@" + confs.Host + ":" + confs.Port + "/")
	if err != nil {
		return rabbit, err
	}
	rabbit.Channel, err = conn.Channel()
	if err != nil {
		return rabbit, err
	}
	return rabbit, err
}

func (s *SynchroService) StartConsume() error {
	q, err := s.rabbit.Channel.QueueDeclare("receipts", true, false, false, false, nil)
	if err != nil {
		return err
	}
	msgs, err := s.rabbit.Channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func () {
		for s.repo.CheckDBConnection() {
			for msg := range msgs {
				var ur templates.UserReceiptMapJSON
				data := msg.Body
				if err := json.Unmarshal(data, &ur); err != nil {
					log.Println("incorrect json urmap: " + err.Error())
					continue
				}
				if err := s.repo.InsertReceipt(ur); err != nil {
					log.Printf("failed to insert receipt: " + err.Error())
					continue
				}
				log.Println("nothing to consume")
			}
			time.Sleep(1 * time.Second)
		}
		log.Fatal("connection with database lost")
	}()
	return nil
}
/*
func (s *SynchroService) processQueue(msgs *<-chan amqp.Delivery) {
	for s.repo.CheckDBConnection() {
		log.Println("xxx")
		for msg := range *msgs {
			var ur templates.UserReceiptMapJSON
			data := msg.Body
			if err := json.Unmarshal(data, &ur); err != nil {
				log.Println("incorrect json urmap: " + err.Error())
				continue
			}
			if err := s.repo.InsertReceipt(ur); err != nil {
				log.Printf("failed to insert receipt: " + err.Error())
				continue
			}
		}
		time.Sleep(1 * time.Second)
	}
	log.Fatal("connection with database lost...")
}*/
