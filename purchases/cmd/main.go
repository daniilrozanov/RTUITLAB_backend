package main

import (
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	templates "purchases/pkg"
	"purchases/pkg/buisness"
	"purchases/pkg/handlers"
	"purchases/pkg/repository"
	"time"
)

func main() {
	time.Sleep(10 * time.Second)

	if err := initConfigs(); err != nil {
		log.Fatalf("Error occured while reading config: %s", err)
	}
	db, err := repository.InitPostgresDB(repository.PostgresConfig{
		Host:     viper.GetString("db_pg.host"),
		Port:     viper.GetString("db_pg.port"),
		Username: viper.GetString("db_pg.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db_pg.dbname"),
		SSLMode:  viper.GetString("db_pg.sslmode"),
	})
	if err != nil {
		log.Fatalf("Error occured while connecting with database: %s", err)
	}
	rabb, err := buisness.NewRabbitStruct(&buisness.RabbitConnectionConfig{
		Host:     viper.GetString("rabbitmq.host"),
		Port:     viper.GetString("rabbitmq.port"),
		Username: viper.GetString("rabbitmq.username"),
		Password: viper.GetString("rabbitmq.password"),
	})
	if err != nil {
		log.Fatalf("Error occured while creating rabbitmq connection: %s", err)
	}

	repo := repository.InitRepositoryLayer(db)
	buis := buisness.InitBuisnessLayer(repo, &rabb)
	handl := handlers.InitHandlersLayer(buis)
	serv := new(templates.Server)

	go buis.StartConsume()
	log.Println("zzz")
	if err := serv.Start(viper.GetString("port"), handl.InitRouting()); err != nil {
		log.Fatalf("Error occured while server tried to start: %s", err.Error())
	}

}

func initConfigs() error{
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}