package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"shops/pkg"
	"shops/pkg/handler"
	_ "shops/pkg/repository"
	"shops/pkg/service"
)

func main(){
	if err := initConfigs(); err != nil {
		log.Fatalf("Error occured while reading of config: %s", err)
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error occured while reading env variables: %s", err)
	}

	/*db, err := repository.InitPostgresDB(repository.PostgresConfig{
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
	*/

	auth := service.NewAuthService(&service.UserServiceConfig{
		Host: viper.GetString("users_service.host"),
		Port: viper.GetString("users_service.port"),
		URN: viper.GetString("users_service.urn"),
		Scheme: viper.GetString("users_service.scheme"),
	})

	//repo := repository.InitRepositoryLayer(db)
	service := service.InitNewService(auth)
	handl := handler.InitNewHandler(service)
	serv := new(pkg.Server)

	if err := serv.Start(viper.GetString("port"), handl.InitRoutes()); err != nil {
		log.Fatalf("Error occured while server tried to start: %s", err.Error())
	}
}

func initConfigs() error{
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}