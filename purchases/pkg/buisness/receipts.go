package buisness

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	templates "purchases/pkg"
	"purchases/pkg/repository"
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
	go func() {
		log.Println("start synchronization")
		for msg := range msgs {
			var ur templates.UserReceiptMapJSON
			data := msg.Body
			if err := json.Unmarshal(data, &ur); err != nil {
				s.throwError("incorrect json urmap: " + err.Error(), q.Name, data)
				return
			}
			if err := s.repo.InsertReceipt(ur); err != nil {
				s.throwError("failed to insert receipt: " + err.Error(), q.Name, data)
				return
			}
			log.Println("receipt synchronized")
		}
		log.Println("stop synchronization")
	}()
	return nil
}

func (s *SynchroService) throwError(errstr, qname string, data []byte) error{
	log.Println(errstr)
	err := s.rabbit.Channel.Publish("", qname, false, false, amqp.Publishing{ContentType: "application/json", Body: data})
	log.Println("synchronization with shops stopped")
	if err != nil {
		return err
	}
	return nil
}

func (s *SynchroService) GetReceipts (userID int) ([]templates.ReceiptJSON, error) {
	return s.repo.GetReceipts(userID)
}