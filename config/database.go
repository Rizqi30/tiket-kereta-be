package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// change detail Postgres local
	dsn := "host=localhost user=postgres password=admin dbname=tiket_kereta port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed connect to db PostgreSQL: ", err)
	}
	log.Println("Success connect to db PostgreSQL")
}