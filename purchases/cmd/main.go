package main

import (
	"github.com/joho/godotenv"
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
	time.Sleep(2 * time.Second)

	if err := initConfigs(); err != nil {
		log.Fatalf("Error occured while reading of config: %s", err)
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error occured while reading env variables: %s", err)
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

	repo := repository.InitRepositoryLayer(db)
	buis := buisness.InitBuisnessLayer(repo)
	handl := handlers.InitHandlersLayer(buis)
	serv := new(templates.Server)
/*
	fmt.Println("Starting the application...")
	ciphertext, err := templates.Encrypt([]byte("daniil"), "LH21tjjg&")
	fmt.Printf("Encrypted: %x\n", ciphertext)
	plaintext, err := templates.Decrypt(ciphertext, "LH21tjjg&")
	fmt.Printf("Decrypted: %s\n", plaintext)
*/
	if err := serv.Start(viper.GetString("port"), handl.InitRouting()); err != nil {
		log.Fatalf("Error occured while server tried to start: %s", err.Error())
	}

}

func initConfigs() error{
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}