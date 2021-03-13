package main

import (
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"shops/pkg"
	"shops/pkg/handler"
	"shops/pkg/repository"
	"shops/pkg/service"
	"time"
)

func main(){
	time.Sleep(11 * time.Second)

	if err := initConfigs(); err != nil {
		log.Fatalf("Error occured while reading of config: %s", err)
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
		log.Fatalf("Error occured while connecting with database: %s", err)
	}
	uConfs := &service.UserServiceConfig{
		Host: viper.GetString("users_service.host"),
		Port: viper.GetString("users_service.port"),
		ConfirmURN: viper.GetString("users_service.confirm_urn"),
		SynchroURN: viper.GetString("users_service.confirm_urn"),
		Scheme: viper.GetString("users_service.scheme"),
	}

	rabb, err := service.NewRabbitStruct(&service.RabbitConnectionConfig{
		Host:     viper.GetString("rabbitmq.host"),
		Port:     viper.GetString("rabbitmq.port"),
		Username: viper.GetString("rabbitmq.username"),
		Password: viper.GetString("rabbitmq.password"),
	})

	defer rabb.Channel.Close()

	if err != nil {
		log.Fatalf("Error occured while connecting with rabbitmq: %s", err)
	}

	repo := repository.NewRepository(db)
	service := service.InitNewService(uConfs, repo, &rabb)
	handl := handler.InitNewHandler(service)
	serv := new(pkg.Server)

	if err := service.SendUnsyncReceiptsToRabbit(); err != nil {
		log.Println("synchronization with purchases failed: ", err.Error())
	}
	if err := service.StartConsume(); err != nil {
		log.Println("synchronization with fabric failed: ", err.Error())
	}

	if err := serv.Start(viper.GetString("port"), handl.InitRoutes()); err != nil {
		log.Fatalf("Error occured while server tried to start: %s", err.Error())
	}
}

func initConfigs() error{
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}