package main

import (
	"fabric/pkg/repository"
	"fabric/pkg/service"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

func main(){
	time.Sleep(11 * time.Second)

	if err := initConfigs(); err != nil {
		log.Fatalf("Error occured while reading of config: %s", err)
	}
	var rabb, err = service.NewRabbitStruct(&service.RabbitConnectionConfig{
		Host:     viper.GetString("rabbitmq.host"),
		Port:     viper.GetString("rabbitmq.port"),
		Username: viper.GetString("rabbitmq.username"),
		Password: viper.GetString("rabbitmq.password"),
	})
	if err != nil {
		log.Println("rabbitmq not connected: ", err.Error())
	}
	cfg, err := service.InitConfig("configs/products.json")
	if err != nil {
		log.Fatal("error while reading configs: ", err.Error())
	}
	db, err := repository.InitPostgresDB(&repository.PostgresConfig{
		Host:     viper.GetString("db_pg.host"),
		Port:     viper.GetString("db_pg.port"),
		Username: viper.GetString("db_pg.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db_pg.dbname"),
		SSLMode:  viper.GetString("db_pg.sslmode"),
	})
	if err != nil {
		log.Fatal("error while connecting with database: ", err.Error())
	}
	repo := repository.NewFabric(db)
	fabric := service.NewFabric(repo, &rabb, cfg)

	log.Println(*cfg)
	fabric.StartProducing()
}


func initConfigs() error{
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}