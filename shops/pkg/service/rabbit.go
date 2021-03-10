package service

import (
	"github.com/streadway/amqp"
	"log"
)

type RabbitConnectionConfig struct {
	Host string
	Port string
	Username string
	Password string
}

func NewRabbitStruct(confs *RabbitConnectionConfig) (RabbitStruct, error) {
	var rabbit RabbitStruct

	log.Println("amqp://"+confs.Username+":"+confs.Password+"@"+confs.Host+":"+confs.Port+"/")
	conn, err := amqp.Dial("amqp://"+confs.Username+":"+confs.Password+"@"+confs.Host+":"+confs.Port+"/")
	if err != nil {
		return rabbit, err
	}
	rabbit.Channel, err = conn.Channel()
	if err != nil {
		return rabbit, err
	}
	return rabbit, nil
}